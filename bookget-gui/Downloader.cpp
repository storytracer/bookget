#include "Downloader.h"
#include "Util.h"
#include "Config.h"
#include "SharedMemory.h"
#include "HttpClient.h"
#include "CheckFailure.h"

void Downloader::Start(HWND hWnd) {
    m_downloaderThread = std::thread([this, hWnd]() {
        // Sleep for 3 seconds, wait for first tab to initialize successfully
        std::this_thread::sleep_for(std::chrono::seconds(3));
        {
            std::lock_guard<std::mutex> lock(m_downloadCounterMutex);

            auto& conf = Config::GetInstance();
            m_downloadCounter = 0;
            m_downloaderMode = conf.GetDownloaderMode();
            m_downloadDir = Util::Utf8ToWide(conf.GetDownloadDir());
            m_sleepTime = conf.GetSleepTime();
            m_maxDownloads = conf.GetMaxDownloads();

            // Load URL list
            if(m_downloaderMode == 0) {
                std::wstring sUrlsFilePath = Util::GetCurrentExeDirectory() + L"\\urls.txt";
                LoadImageUrlsFromFile(sUrlsFilePath);
            }
        }

  

        if (!m_targetUrls.empty())
        {
            // Start downloading
            DownloadNextImage(hWnd);
        }

        m_workerThreadId = GetCurrentThreadId(); // Save thread ID
            
        // Message loop (required for receiving PostThreadMessage)
         while (!m_stopThread.load(std::memory_order_relaxed)) {
            MSG msg;
            // Set 100ms timeout, avoid long-term blocking
            if (PeekMessage(&msg, nullptr, 0, 0, PM_REMOVE)) {
                if (msg.message == WM_DOWNLOAD_URL) {
                    std::wstring* pUrl = reinterpret_cast<std::wstring*>(msg.lParam);
                    //DownloadFile(*pUrl, *pUrl);
                    delete pUrl;
                }
            }
            std::this_thread::sleep_for(std::chrono::milliseconds(100));
        }
    });
}
Downloader::~Downloader(){
    Stop();
}

void Downloader::Stop() {
    m_stopThread = true;
    if (m_downloaderThread.joinable()) {
        m_downloaderThread.join();
    }
    m_stopThread = false;
}

void Downloader::RequestDownload(const std::wstring& url) {
    // Deep copy data and send to worker thread
    std::wstring* pUrl = new std::wstring(url);
    ::PostThreadMessage(m_workerThreadId, WM_DOWNLOAD_URL, 0, reinterpret_cast<LPARAM>(pUrl));

}
bool Downloader::ShouldInterceptRequest(const std::wstring& sUrl){
    
    // Skip local paths (file://, http://localhost, 127.0.0.1, etc.)
    if (Util::IsLocalUri(sUrl)) {
        return false; // Do not handle local URIs
    }

    // 1. Check if URL matches target URL list
    bool urlMatch = false;
    for (const auto& targetUrl : m_targetUrls)
    {
        if (sUrl.find(targetUrl) != std::wstring::npos)
        {
            urlMatch = true;
            break;
        }
    }
    // 2. Check if URL matches config.yaml URL
    auto& conf = Config::GetInstance();
    std::string narrow_url = Util::WideToUtf8(sUrl);
    for (const auto& site : conf.GetSiteConfigs()) {
        // URL before HTTP request
        if (site.intercept == 0 && Util::matchUrlPattern(site.url, narrow_url) ) {
             urlMatch = true;
            break;
        }
    }

    // 3. Check URL extension
    //for (const auto& ext : m_targetExtensions) {
    //    if (sUrl.size() >= ext.size() && 
    //        _wcsicmp(sUrl.substr(sUrl.size() - ext.size()).c_str(), ext.c_str()) == 0) {
    //        urlMatch = true;
    //        break;
    //    }
    //}
    
    return urlMatch;
}

bool Downloader::ShouldInterceptResponse(const std::wstring& sUrl)
{
    // Skip local paths (file://, http://localhost, 127.0.0.1, etc.)
    if (Util::IsLocalUri(sUrl)) {
        return false; // Do not handle local URIs
    }

    // 1. Check if URL matches target URL list
    bool urlMatch = false;
    for (const auto& targetUrl : m_targetUrls)
    {
        if (sUrl.find(targetUrl) != std::wstring::npos)
        {
            urlMatch = true;
            break;
        }
    }
    // 2. Check if URL matches config.yaml URL
    auto& conf = Config::GetInstance();
    std::string narrow_url = Util::WideToUtf8(sUrl);
    for (const auto& site : conf.GetSiteConfigs()) {
        // URL after HTTP request
        if (site.intercept == 1 && Util::matchUrlPattern(site.url, narrow_url) ) {
             urlMatch = true;
            break;
        }
    }

    return urlMatch;
}

bool Downloader::ShouldInterceptContentType(const std::wstring& contentType, const std::wstring& contentLength)
{
    bool isCanDownload = false;
    // Check Content-Type
    if(!contentType.empty()) {
        for (const auto& ext : m_targetContentTypes) {
            if (contentType.size() >= ext.size() && 
                contentType.find(ext) != std::wstring::npos) {
                isCanDownload = true;
                break;
            }
        }
    }

    // Check Content-Length
    if(!contentLength.empty()) {
        ULONGLONG length = 0;
        if (swscanf_s(contentLength.c_str(), L"%llu", &length) == 1) {
            // Set reasonable image size range (10KB - 20MB)
            isCanDownload = (length >= 10240 && length <= 20 * 1024 * 1024);
        }
    }

    return isCanDownload;
}


std::wstring Downloader::GetFilePath(const std::wstring& sUrl)
{
    std::wstring filePath;

    bool isSharedDataURL = false;
    // Read shared memory
    auto* sharedData = SharedMemory::GetInstance().GetMutex();
    if (sharedData != nullptr) {
        isSharedDataURL = sharedData->ImageReady && sharedData->ImagePath && sharedData->PID != GetCurrentProcessId();
        filePath.assign(sharedData->ImagePath);
        if (isSharedDataURL)
             m_downloaderMode = 2;
        SharedMemory::GetInstance().ReleaseMutex();
     }

    if (m_downloaderMode == 0 || m_downloaderMode == 1) {
         // Get next sequence number
        std::lock_guard<std::mutex> lock(m_downloadCounterMutex);
        int currentCount = ++m_downloadCounter;
 

        // Default extension
        auto& conf = Config::GetInstance();
        std::wstring ext = Util::Utf8ToWide(conf.GetDefaultExt());
        bool useDefaultExt = false;
        std::string narrow_url = Util::WideToUtf8(sUrl);
        for (const auto& site : conf.GetSiteConfigs()) {
            if (Util::matchUrlPattern(site.url, narrow_url) ) {
               ext = Util::Utf8ToWide(site.ext);
               useDefaultExt = true;
               break;
            }
        }

        std::wstringstream filename;
        filename << m_downloadDir << L"\\"  
            << std::setw(4) << std::setfill(L'0') << currentCount;

        if(!useDefaultExt) {
            // Try to get file extension from URL
            size_t dotPos = sUrl.find_last_of(L'.');
            if (dotPos != std::wstring::npos)
            {
                std::wstring ext_ = sUrl.substr(dotPos);
                if (ext.length() <= 5) 
                {
                    filename << ext_;
                }
                else {
                    filename << ext;
                }
            }
        }
        else {
            filename << ext;
        }
        filePath.assign(filename.str());
    }
 
    if (m_downloadCounter >=  m_maxDownloads)
    {
        OutputDebugString(L"Exceeded max_downloads limit set in config.yaml\n");
        return L"";
    }
    return filePath;
}


// 2. Load URLs from file
void Downloader::LoadImageUrlsFromFile(const std::wstring& sUrlsFilePath)
{
    std::wifstream file;
    if (sUrlsFilePath.empty())
        return;

    file.open(sUrlsFilePath);
    if (!file.is_open())
    {
        OutputDebugString(L"Error: Could not open any urls file (global or local)\n");
        return;
    }

    m_downloaderMode = 0;

    m_targetUrls.clear();
    std::wstring line;
    while (std::getline(file, line))
    {
        if (!line.empty())
        {
            m_targetUrls.emplace_back(line);
        }
    }
}

// 3. Download next page
void Downloader::DownloadNextImage(HWND hWnd)
{
    int currentIndex = m_downloadCounter;

    if (currentIndex >= m_targetUrls.size() || currentIndex >=  m_maxDownloads)
    {
        OutputDebugString(L"All downloads completed\n");
        return;
    }
   

    try {
        std::unique_ptr<std::wstring> pUrl = std::make_unique<std::wstring>(m_targetUrls.at(currentIndex));
        ::PostMessage(
            hWnd,
            WM_DOWNLOAD_URL,
            0,
            reinterpret_cast<LPARAM>(pUrl.release()) // Transfer ownership
        );
    } catch (const std::out_of_range&) {
        //::PostMessage(m_hWnd, WM_ERR, 0, (LPARAM)L"Index out of range");
    }
    
}



bool Downloader::DownloadFile(const wchar_t* url, ICoreWebView2HttpRequestHeaders* headers)
{
    std::wstring filePath = GetFilePath(url);
    std::vector<std::pair<std::string, std::string>> headersVec = {};

    wil::com_ptr<ICoreWebView2HttpHeadersCollectionIterator> iterator;
    if (SUCCEEDED(headers->GetIterator(&iterator))) {
        BOOL hasCurrent = FALSE;
        while (SUCCEEDED(iterator->get_HasCurrentHeader(&hasCurrent)) && hasCurrent) {
            wil::unique_cotaskmem_string name, value;
            if (SUCCEEDED(iterator->GetCurrentHeader(&name, &value))) {
                headersVec.emplace_back(Util::WideToUtf8(name.get()), Util::WideToUtf8(value.get()));
            }
            iterator->MoveNext(&hasCurrent);
        }
    }

    try {
        asio::io_context io_context;
        ssl::context ssl_ctx(ssl::context::tls_client);
        
        ssl_ctx.set_verify_mode(ssl::verify_none);
        
        HttpClient httpClient(io_context, ssl_ctx);
        
        std::string sUrl_u8 = Util::WideToUtf8(url);
        std::string filePath_u8 = Util::WideToUtf8(filePath);
        if (httpClient.download(sUrl_u8, filePath_u8, headersVec)) {
            OutputDebugString(L"Download completed successfully!");
            return true;
        } else {
            OutputDebugString(L"Download failed");
            return false;
        }
    } catch (std::exception& e) {
            Util::DebugPrintException(e);
        return false;
    }
}
