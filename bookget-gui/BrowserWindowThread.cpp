#include "BrowserWindow.h"
#include "SharedMemory.h"

void BrowserWindow::StartBackgroundThread() {
    std::lock_guard<std::mutex> lock(m_threadMutex);
    StopBackgroundThread();  // Ensure previous thread is stopped

    m_sharedMemoryThread = std::thread([this]() {
        while (!m_stopThread.load(std::memory_order_relaxed)) {
            // Read simplified shared memory data
            auto data = SharedMemory::GetInstance().Read();
            
             // Check if there are new image paths to process
            if (data.URLReady && data.ImageReady && data.PID != GetCurrentProcessId()) {
                auto* sharedData = new SharedMemoryDataMini();
                sharedData->ImageReady = data.ImageReady;
                sharedData->URLReady = data.URLReady;
                sharedData->PID = data.PID;
                wcsncpy_s(sharedData->ImagePath, _countof(sharedData->ImagePath), data.ImagePath, _TRUNCATE);
                wcsncpy_s(sharedData->URL, _countof(sharedData->URL), data.URL, _TRUNCATE);
                
                PostMessage(m_hWnd, WM_APP_UPDATE_UI, 0, reinterpret_cast<LPARAM>(sharedData));
            }
            // Check if there are new URLs to process
            else if (data.URLReady && data.PID != GetCurrentProcessId()) {
                // Only pass necessary data to UI thread
                auto* sharedData = new SharedMemoryDataMini();
                sharedData->URLReady = data.URLReady;
                sharedData->PID = data.PID;
                wcsncpy_s(sharedData->URL, _countof(sharedData->URL), data.URL, _TRUNCATE);
                
                PostMessage(m_hWnd, WM_APP_UPDATE_UI, 0, reinterpret_cast<LPARAM>(sharedData));
            }
                   
            // Reduce CPU usage
            std::this_thread::sleep_for(std::chrono::milliseconds(200));
        }
    });

    m_downloader.Start(m_hWnd);
}

void BrowserWindow::StopBackgroundThread() {
    m_stopThread = true; // Set stop flag first
    // Ensure thread is completely stopped
    if (m_sharedMemoryThread.joinable()) {
        m_sharedMemoryThread.join();
    }

    m_downloader.Stop();

    m_stopThread = false;  // Reset flag
}



void BrowserWindow::HandleSharedMemoryUpdate(LPARAM lParam) {
    // Get the passed data
    auto* data = reinterpret_cast<SharedMemoryData*>(lParam);

    auto* sharedData = SharedMemory::GetInstance().GetMutex();
    if (sharedData == nullptr) {
        delete data;
        return;
     }
    
    // Handle image download mode
    if (sharedData->URLReady && sharedData->ImageReady && m_tabs.find(m_activeTabId) != m_tabs.end() && sharedData->PID != GetCurrentProcessId()) {
        // Reset flags
        sharedData->URLReady = false;
        m_downloader.Reset(sharedData->URL, 2);
        // Configure download handler
        m_tabs.at(m_activeTabId)->SetupWebViewListeners();
        // Navigate to URL
        m_tabs.at(m_activeTabId)->m_contentWebView->Navigate(sharedData->URL);
    } 
    // Handle normal URL navigation
    else if (sharedData->URLReady && m_tabs.find(m_activeTabId) != m_tabs.end()) {
        // Reset flags
        sharedData->URLReady = false;
            
        // Navigate to URL
        m_tabs.at(m_activeTabId)->m_contentWebView->Navigate(sharedData->URL);
    }
        
    
    // Clean up memory
    delete data;
    SharedMemory::GetInstance().ReleaseMutex();

}

