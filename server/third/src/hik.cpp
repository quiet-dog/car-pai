#include "hik.h"
#include "HCNetSDK.h"
#include <string.h>
#include <cstring>
class HikvisionCamera : public CameraDevice
{
private:
    int lUserID;

public:
    explicit HikvisionCamera(const DeviceConfig &config) : CameraDevice(config)
    {
    }

    // nvc登陆
    bool login()
    {
        NET_DVR_USER_LOGIN_INFO struLoginInfo = {0};
        NET_DVR_DEVICEINFO_V40 struDeviceInfoV40 = {0};
        struLoginInfo.bUseAsynLogin = false;

        // 端口转换并做错误处理
        try
        {
            struLoginInfo.wPort = static_cast<WORD>(std::stoi(config.port));
        }
        catch (const std::exception &e)
        {
            printf("Error: Invalid port format: %s\n", e.what());
            this->config.online = false;
            return;
        }

        // 复制 IP 地址并确保字符串终止
        std::strncpy(struLoginInfo.sDeviceAddress, config.ip.c_str(), NET_DVR_DEV_ADDRESS_MAX_LEN - 1);
        struLoginInfo.sDeviceAddress[NET_DVR_DEV_ADDRESS_MAX_LEN - 1] = '\0';

        // 复制用户名和密码（假设它们存储在 config.user 和 config.password 中）
        std::strncpy(struLoginInfo.sUserName, config.username.c_str(), NAME_LEN - 1);
        struLoginInfo.sUserName[NAME_LEN - 1] = '\0';

        std::strncpy(struLoginInfo.sPassword, config.password.c_str(), NAME_LEN - 1);
        struLoginInfo.sPassword[NAME_LEN - 1] = '\0';

        // 执行登录
        int lUserID = NET_DVR_Login_V40(&struLoginInfo, &struDeviceInfoV40);
        if (lUserID < 0)
        {
            printf("pyd---Login error, %d\n", NET_DVR_GetLastError());
            this->config.online = false;

            // 仅当已初始化时才调用 Cleanup，避免二次清理
            NET_DVR_Cleanup();
            return false;
        }
        this->lUserID = lUserID;
        return true;
    }

    // 布防
    bool fortify()
    {
        // 假设登录成功
        if (this->lUserID < 0)
        {
            printf("pyd---Login error, %d\n", NET_DVR_GetLastError());
            return false;
        }

        return true;
    }

    void connect() override
    {
    }

    void disconnect() override
    {
    }

    void captureImage() override
    {
    }

    void streamVideo() override
    {
    }
};