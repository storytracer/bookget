#pragma once

#include <yaml-cpp/yaml.h>
#include <memory>

class Config
{
public:
    // Delete copy constructor and assignment operator
    Config(const Config&) = delete;
    Config& operator=(const Config&) = delete;
    
    // Get singleton instance
    static Config& GetInstance() {
        static Config instance;
        return instance;
    }

    // Configuration management
    struct SiteConfig {
        std::string url;
        std::string script;
        int intercept; 
        std::string ext;
        std::string description;
        int downloaderMode;
    };

    // Public interface
    bool Load(const std::string& configPath);
    std::string GetDownloadDir();
    std::string GetDefaultExt();
    int GetMaxDownloads();
    int GetSleepTime();
    int GetDownloaderMode();
    const std::vector<SiteConfig>& GetSiteConfigs();

private:
    // PIMPL 实现
    struct ConfigImpl {
        // Global settings
        std::string downloadDir = "downloads";
        int maxDownloads = 1000;
        int sleepTime = 3;
        int downloaderMode = 1;    // Download mode 0=urls.txt | 1=auto listen images | 2=shared memory URL
        std::string fileExt = ".jpg";
     
        std::vector<SiteConfig> siteConfigs;

        // Load YAML configuration file
        bool Load(const std::string& configPath);
    };

    Config();  // Private constructor
    ~Config() = default;
    
    // PIMPL implementation
    std::unique_ptr<ConfigImpl> pImpl;
};
