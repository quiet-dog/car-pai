#ifndef HIK_H
#define HIK_H

#include "device.h"

class HikvisionCamera : public CameraDevice
{
public:
    explicit HikvisionCamera(const DeviceConfig &config) : CameraDevice(config) {}

    void connect() override;

    void disconnect() override;

    void captureImage() override;

    void streamVideo() override;
};

#endif // HIKVISION_CAMERA_H