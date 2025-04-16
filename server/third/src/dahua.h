#ifndef DAHUA_H
#define DAHUA_H

#include "device.h"

class DahuaCamera : public CameraDevice
{
public:
    explicit DahuaCamera(const DeviceConfig &config) : CameraDevice(config) {}

    void connect() override
    {
        // 大华设备的连接实现
    }

    void disconnect() override
    {
        // 大华设备的断开实现
    }

    void captureImage() override
    {
        // 大华设备的拍照实现
    }

    void streamVideo() override
    {
        // 大华设备的视频流实现
    }
};

#endif // DAHUA_CAMERA_H
