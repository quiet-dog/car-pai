#include "device.h"

class CameraDevice
{
protected:
    DeviceConfig config;

public:
    explicit CameraDevice(const DeviceConfig &config) : config(config) {}

    virtual ~CameraDevice() = default;
    // 定义一些设备通用操作
    virtual void connect() = 0;
};