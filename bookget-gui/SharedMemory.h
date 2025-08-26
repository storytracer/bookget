#pragma once
#include <mutex>
#include <windows.h>

// Shared memory structure
#include <cstdint>
#include <string>
#include <tchar.h>

#pragma pack(push, 1)  // Ensure no padding bytes
struct SharedMemoryData {
	uint32_t URLReady;
	uint32_t HTMLReady;
	uint32_t CookiesReady;
	uint32_t ImageReady;
	uint32_t PID;
	wchar_t URL[1024];                // Fixed size buffer
	wchar_t ImagePath[1024];          // MAX_PATH for image path
	wchar_t Cookies[4096];            // 4KB for cookies
	wchar_t HTML[1024 * 1024 * 8];    // 8MB for HTML
};
struct SharedMemoryDataMini {
	uint32_t URLReady;
	uint32_t HTMLReady;
	uint32_t CookiesReady;
	uint32_t ImageReady;
	uint32_t PID;
	wchar_t URL[1024];                // Fixed size buffer
	wchar_t ImagePath[1024];          // MAX_PATH for image path
};
#pragma pack(pop)  // Restore default alignment


class SharedMemory
{

public:
    // Get singleton instance
    static SharedMemory& GetInstance() {
        static SharedMemory instance;
        return instance;
    }


     bool Init();
     void Cleanup();
     void WriteHtml(const std::wstring& html);
     void WriteCookies(const std::wstring& cookies);
     void WriteImagePath(const std::wstring& imagePath);
     SharedMemoryDataMini Read();
     SharedMemoryData* Get();
     SharedMemoryData* GetMutex();
     void ReleaseMutex();


private:

     HANDLE m_hSharedMemory;        
     LPVOID m_pSharedMemory;          
     HANDLE m_hSharedMemoryMutex;      

    // Optimization: use constexpr instead of static const
     const wchar_t* m_sharedMemoryName = L"Local\\WebView2SharedMemory";
     const wchar_t* m_sharedMemoryMutexName = L"Local\\WebView2SharedMemoryMutex";
     DWORD m_sharedMemorySize = sizeof(SharedMemoryData);
};

