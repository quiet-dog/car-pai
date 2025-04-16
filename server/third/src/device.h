#ifndef DEVICE_H
#define DEVICE_H

#include <string>

struct DeviceConfig
{
    std::string ip;
    std::string port;
    std::string username;
    std::string password;
    bool online = false;

    DeviceConfig(const std::string &ip, const std::string &port, const std::string &username, const std::string &password)
        : ip(ip), port(port), username(username), password(password) {}
};

class CameraDevice
{

protected:
    DeviceConfig config;

public:
    explicit CameraDevice(const DeviceConfig &config) : config(config) {}

    virtual ~CameraDevice() = default;
    // 定义一些设备通用操作
    virtual void connect() = 0;

    virtual void disconnect() = 0;

    virtual void captureImage() = 0;

    virtual void streamVideo() = 0;
};

#endif // DEVICE_H