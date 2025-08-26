// Copyright (C) Microsoft Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#pragma once

#include "framework.h"
#include "Tab.h"
#include <chrono>
#include <thread>
#include "SharedMemory.h"
#include <mutex>
#include <atomic>
#include <memory>
#include <vector>
#include <string>
#include <wrl.h>
#include <wil/com.h>
#include <cpprest/json.h>
#include "Util.h"

#include <yaml-cpp/yaml.h>
#include <filesystem>

#include "shlobj.h"
#include <Urlmon.h>
#pragma comment (lib, "Urlmon.lib")
#include "env.h"

#include <shlwapi.h> // for PathCombine
#include "CheckFailure.h"
#pragma comment(lib, "shlwapi.lib")

#include "Downloader.h"


using namespace Microsoft::WRL;
using Microsoft::WRL::Callback;


class BrowserWindow
{
public:
    BrowserWindow(){};
    ~BrowserWindow();

    // Window size constants
    static const int c_uiBarHeight = 70;
    static const int c_optionsDropdownHeight = 108;
    static const int c_optionsDropdownWidth = 200;

    // Window management
    static ATOM RegisterClass(_In_ HINSTANCE hInstance);
    static LRESULT CALLBACK WndProcStatic(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam);
    LRESULT CALLBACK WndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam);
    static BOOL LaunchWindow(_In_ HINSTANCE hInstance, _In_ int nCmdShow);



    // WebView management
    HRESULT OpenWindowTab(wchar_t *webUrl, bool isTab = false);
    HRESULT ExecuteScriptFile(const std::wstring& scriptPath, ICoreWebView2* webview);

    // Tab management
    HRESULT HandleTabURIUpdate(size_t tabId, ICoreWebView2* webview);
    HRESULT HandleTabHistoryUpdate(size_t tabId, ICoreWebView2* webview);
    HRESULT HandleTabNavStarting(size_t tabId, ICoreWebView2* webview);
    HRESULT HandleTabNavCompleted(size_t tabId, ICoreWebView2* webview, ICoreWebView2NavigationCompletedEventArgs* args);
    HRESULT HandleTabSecurityUpdate(size_t tabId, ICoreWebView2* webview, ICoreWebView2DevToolsProtocolEventReceivedEventArgs* args);
    void HandleTabCreated(size_t tabId, bool shouldBeActive);
    HRESULT HandleTabMessageReceived(size_t tabId, ICoreWebView2* webview, ICoreWebView2WebMessageReceivedEventArgs* eventArgs);
    void HandleSharedMemoryUpdate(LPARAM lParam);


    // Worker thread
    void StartBackgroundThread();
    void StopBackgroundThread();

    // Utility tools
    static std::wstring GetAppDataDirectory();
    std::wstring GetFullPathFor(LPCWSTR relativePath);
    std::wstring GetFilePathAsURI(std::wstring fullPath);
    static std::wstring GetUserDataDirectory();
    int GetDPIAwareBound(int bound);
    static void CheckFailure(HRESULT hr, LPCWSTR errorMessage = L"");

private:
    // Worker thread
    std::mutex m_threadMutex;
    std::thread m_sharedMemoryThread;
    std::atomic<bool> m_stopThread{false};

protected:
    // Window resources
    HINSTANCE m_hInst = nullptr;
    HWND m_hWnd = nullptr;
    static WCHAR s_windowClass[MAX_LOADSTRING];
    static WCHAR s_title[MAX_LOADSTRING];
    int m_minWindowWidth = 0;
    int m_minWindowHeight = 0;

public:
    // WebView resources
    wil::com_ptr<ICoreWebView2Environment> m_uiEnv;
    wil::com_ptr<ICoreWebView2Environment> m_contentEnv;
protected:
    wil::com_ptr<ICoreWebView2Controller> m_controlsController;
    wil::com_ptr<ICoreWebView2Controller> m_optionsController;
    wil::com_ptr<ICoreWebView2> m_controlsWebView;
    wil::com_ptr<ICoreWebView2> m_optionsWebView;
    std::map<size_t, std::unique_ptr<Tab>> m_tabs;
    size_t m_activeTabId = 0;

    // WebView event handling
    EventRegistrationToken m_controlsUIMessageBrokerToken;
    EventRegistrationToken m_optionsUIMessageBrokerToken;
    EventRegistrationToken m_controlsZoomToken;
    EventRegistrationToken m_optionsZoomToken;
    EventRegistrationToken m_lostOptionsFocus;
    EventRegistrationToken m_newWindowRequestedToken;
    wil::com_ptr<ICoreWebView2WebMessageReceivedEventHandler> m_uiMessageBroker;


public:
    // Initialization methods
    BOOL InitInstance(HINSTANCE hInstance, int nCmdShow);
    HRESULT InitUIWebViews();
    HRESULT CreateBrowserControlsWebView();
    HRESULT CreateBrowserOptionsWebView();
    void SetUIMessageBroker();
    HRESULT ResizeUIWebViews();
    void UpdateMinWindowSize();
    HRESULT PostJsonToWebView(web::json::value jsonObj, ICoreWebView2* webview);
    HRESULT SwitchToTab(size_t tabId);


    // Cache and Cookie management
    HRESULT ClearContentCache();
    HRESULT ClearControlsCache();
    HRESULT ClearContentCookies();
    HRESULT ClearControlsCookies();

    
    // Download management
    Downloader m_downloader;

    // Download handling
    HRESULT HandleTabWebResourceResponseReceived(ICoreWebView2* sender, 
        ICoreWebView2WebResourceResponseReceivedEventArgs* args);
    bool ShouldInterceptResponse(const std::wstring& sUrl, 
        ICoreWebView2WebResourceResponseView* response);
    bool DownloadFile(const std::wstring& sUrl, IStream *content);
    bool DownloadFile(const std::wstring& sUrl, ICoreWebView2HttpRequestHeaders *headers);

};
