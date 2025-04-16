#ifndef GATEWAY_H
#define GATEWAY_H

#include "device.h"
#include <map>
#include <memory>
#include <string>
#include <iostream>

class Gateway
{
private:
    std::map<std::string, std::unique_ptr<CameraDevice>> devices;

public:
    // 添加设备
    void addDevice(const std::string &ip, std::unique_ptr<CameraDevice> device)
    {
        devices[ip] = std::move(device);
    }

    // 移除设备
    void removeDevice(const std::string &ip)
    {
        if (devices.erase(ip) == 0)
        {
            std::cerr << "Device with IP " << ip << " not found!" << std::endl;
        }
    }

    // 获取设备
    CameraDevice *getDevice(const std::string &ip)
    {
        auto it = devices.find(ip);
        if (it != devices.end())
        {
            return it->second.get();
        }
        std::cerr << "Device with IP " << ip << " not found!" << std::endl;
        return nullptr;
    }

    // 连接设备
    void connectDevice(const std::string &ip)
    {
        CameraDevice *device = getDevice(ip);
        if (device)
        {
            device->connect();
        }
    }

    // 断开设备
    void disconnectDevice(const std::string &ip)
    {
        CameraDevice *device = getDevice(ip);
        if (device)
        {
            device->disconnect();
        }
    }

    // 拍照
    void captureImage(const std::string &ip)
    {
        CameraDevice *device = getDevice(ip);
        if (device)
        {
            device->captureImage();
        }
    }

    // 视频流
    void streamVideo(const std::string &ip)
    {
        CameraDevice *device = getDevice(ip);
        if (device)
        {
            device->streamVideo();
        }
    }

    // 列出所有设备
    void listDevices()
    {
        std::cout << "Devices in Gateway:\n";
        for (const auto &pair : devices)
        {
            std::cout << " - " << pair.first << std::endl;
        }
    }
};

#endif // GATEWAY_H
