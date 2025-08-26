#include "SharedMemory.h"


SharedMemoryData* SharedMemory::Get(){
    return static_cast<SharedMemoryData*>(m_pSharedMemory);
};

SharedMemoryData* SharedMemory::GetMutex(){
    DWORD waitResult = ::WaitForSingleObject(m_hSharedMemoryMutex, 5000);
    if (waitResult != WAIT_OBJECT_0) {
        OutputDebugString(L"Failed to acquire mutex \n");
        return nullptr;
    }
    return static_cast<SharedMemoryData*>(m_pSharedMemory);
};

void SharedMemory::ReleaseMutex(){
    ::ReleaseMutex(m_hSharedMemoryMutex);
};

// Initialize shared memory and mutex
bool SharedMemory::Init()
{
    // Create mutex
    m_hSharedMemoryMutex = CreateMutexW(nullptr, FALSE, m_sharedMemoryMutexName);
    if (m_hSharedMemoryMutex == nullptr)
    {
        OutputDebugString(L"Failed to create shared memory mutex\n");
        return false;
    }

    // Wait to acquire mutex
    DWORD waitResult = ::WaitForSingleObject(m_hSharedMemoryMutex, 5000); // 5 second timeout
    if (waitResult != WAIT_OBJECT_0)
    {
        OutputDebugString(L"Failed to acquire shared memory mutex\n");
        return false;
    }

    // Create shared memory
    m_hSharedMemory = CreateFileMappingW(
        INVALID_HANDLE_VALUE,
        NULL,
        PAGE_READWRITE,
        0,
        m_sharedMemorySize,
        m_sharedMemoryName);

    if (m_hSharedMemory == nullptr)
    {
        OutputDebugString(L"Failed to create shared memory\n");
        ::ReleaseMutex(m_hSharedMemoryMutex);
        return false;
    }

    // Map shared memory view
    m_pSharedMemory = MapViewOfFile(
        m_hSharedMemory,
        FILE_MAP_ALL_ACCESS,
        0,
        0,
        m_sharedMemorySize);

    if (m_pSharedMemory == nullptr)
    {
        OutputDebugString(L"Failed to map view of shared memory\n");
        CloseHandle(m_hSharedMemory);
        m_hSharedMemory = nullptr;
        ::ReleaseMutex(m_hSharedMemoryMutex);
        return false;
    }

     // Initialize shared memory structure
    SharedMemoryData* sharedData = static_cast<SharedMemoryData*>(m_pSharedMemory);
    ZeroMemory(sharedData, m_sharedMemorySize); // Clear entire structure
    sharedData->PID = GetCurrentProcessId();

  

    // Release mutex
    ::ReleaseMutex(m_hSharedMemoryMutex);

    return true;
}

SharedMemoryDataMini SharedMemory::Read() {
    SharedMemoryDataMini data;
    ZeroMemory(&data, sizeof(data)); // Initialize structure

    // Acquire mutex
    DWORD waitResult = WaitForSingleObject(m_hSharedMemoryMutex, 5000);
    if (waitResult != WAIT_OBJECT_0) {
        OutputDebugString(L"Failed to acquire mutex for reading shared memory\n");
        return data; // Return empty data
    }

    if (m_pSharedMemory) {
        SharedMemoryDataMini* sharedData = static_cast<SharedMemoryDataMini*>(m_pSharedMemory);
        
        // Copy basic flags
        data.URLReady = sharedData->URLReady;
        data.HTMLReady = sharedData->HTMLReady;
        data.CookiesReady = sharedData->CookiesReady;
        data.ImageReady = sharedData->ImageReady;
        data.PID = sharedData->PID;

        // Safely copy strings
        wcsncpy_s(data.URL, _countof(data.URL), sharedData->URL, _TRUNCATE);
        //wcsncpy_s(data.HTML, _countof(data.HTML), sharedData->HTML, _TRUNCATE);
        //wcsncpy_s(data.Cookies, _countof(data.Cookies), sharedData->Cookies, _TRUNCATE);
        wcsncpy_s(data.ImagePath, _countof(data.ImagePath), sharedData->ImagePath, _TRUNCATE);
    }

    ::ReleaseMutex(m_hSharedMemoryMutex);
    return data;
}


// Write HTML to shared memory
void SharedMemory::WriteHtml(const std::wstring& html) {
    if (m_pSharedMemory == nullptr) return;

    DWORD waitResult = WaitForSingleObject(m_hSharedMemoryMutex, 5000);
    if (waitResult != WAIT_OBJECT_0) {
        OutputDebugString(L"Failed to acquire mutex for writing HTML\n");
        return;
    }

    SharedMemoryData* sharedData = static_cast<SharedMemoryData*>(m_pSharedMemory);
    
    // Clear existing data
    ZeroMemory(sharedData->HTML, sizeof(sharedData->HTML));
    
    // Copy new data
    size_t destSize = sizeof(sharedData->HTML) / sizeof(wchar_t);
    size_t copySize = min(html.size(), destSize - 1);
    
    // Use secure copy function
    wcsncpy_s(sharedData->HTML, destSize, html.c_str(), copySize);
    sharedData->HTML[copySize] = L'\0'; // Ensure null termination


    sharedData->HTMLReady = true;
    sharedData->URLReady = false;
    sharedData->PID = GetCurrentProcessId();

    ::ReleaseMutex(m_hSharedMemoryMutex);
}

// Write Cookies to shared memory
void SharedMemory::WriteCookies(const std::wstring& cookies)
{
    if (m_pSharedMemory == nullptr)
        return;

    DWORD waitResult = WaitForSingleObject(m_hSharedMemoryMutex, 5000);
    if (waitResult != WAIT_OBJECT_0)
    {
        OutputDebugString(L"Failed to acquire mutex for writing cookies\n");
        return;
    }

    SharedMemoryData* sharedData = static_cast<SharedMemoryData*>(m_pSharedMemory);
    
    // Clear existing data
    ZeroMemory(sharedData->Cookies, sizeof(sharedData->Cookies));
    
    // Copy new data
    size_t destSize = sizeof(sharedData->Cookies) / sizeof(wchar_t);
    size_t copySize = min(cookies.size(), destSize - 1);
    
    // Use secure copy function
    wcsncpy_s(sharedData->Cookies, destSize, cookies.c_str(), copySize);
    sharedData->Cookies[copySize] = L'\0'; // Ensure null termination

    sharedData->CookiesReady = true;
    sharedData->PID = GetCurrentProcessId(); 

    ::ReleaseMutex(m_hSharedMemoryMutex);
}

// Write IMAGEPATH to shared memory
void SharedMemory::WriteImagePath(const std::wstring& imagePath)
{
    if (m_pSharedMemory == nullptr)
        return;

    // Acquire mutex
    DWORD waitResult = WaitForSingleObject(m_hSharedMemoryMutex, 5000);
    if (waitResult != WAIT_OBJECT_0)
    {
        OutputDebugString(L"Failed to acquire mutex for writing HTML\n");
        return;
    }

    SharedMemoryData* sharedData = static_cast<SharedMemoryData*>(m_pSharedMemory);

      // Clear existing data
    ZeroMemory(sharedData->ImagePath, sizeof(sharedData->ImagePath));
    
    // Copy new data
    size_t destSize = sizeof(sharedData->ImagePath) / sizeof(wchar_t);
    size_t copySize = min(imagePath.size(), destSize - 1);
    
    // Use secure copy function
    wcsncpy_s(sharedData->ImagePath, destSize, imagePath.c_str(), copySize);
    sharedData->ImagePath[copySize] = L'\0'; // Ensure null termination

    sharedData->ImageReady = false;
    sharedData->URLReady = false;
    sharedData->PID = GetCurrentProcessId(); 

    ::ReleaseMutex(m_hSharedMemoryMutex);
}

// Clean up shared memory resources
void SharedMemory::Cleanup()
{
    // Acquire mutex
    if (m_hSharedMemoryMutex)
    {
        WaitForSingleObject(m_hSharedMemoryMutex, INFINITE);
    }

    // Clean up shared memory mapping
    if (m_pSharedMemory)
    {
        // Zero out entire shared memory area
        ZeroMemory(m_pSharedMemory, m_sharedMemorySize);

        UnmapViewOfFile(m_pSharedMemory);
        m_pSharedMemory = nullptr;
    }
    
    // Close shared memory handle
    if (m_hSharedMemory)
    {
        CloseHandle(m_hSharedMemory);
        m_hSharedMemory = nullptr;
    }

    // Release mutex and close handle
    if (m_hSharedMemoryMutex)
    {
        ::ReleaseMutex(m_hSharedMemoryMutex);
        CloseHandle(m_hSharedMemoryMutex);
        m_hSharedMemoryMutex = nullptr;
    }
}
