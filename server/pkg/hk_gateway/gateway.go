package hk_gateway

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type HikGateway struct {
	hikMap sync.Map
	// 注册接受广播的通道
	broadClient sync.Map

	// udp设备列表
	devices []Device
}

type HikInfo struct {
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsConnect   bool   `json:"is_connect"`
	LongConnect bool   `json:"long_connect"`
}

func NewHikGateway() *HikGateway {
	gateway := &HikGateway{}
	// go gateway.GetUDPClient()
	// go gateway.WatchLongConnect()
	return gateway
}

// func (h *HikGateway) WatchLongConnect() {
// 	for {
// 		time.Sleep(time.Minute * 2)
// 		h.hikMap.Range(func(key, value any) bool {
// 			for _, v := range h.devices {
// 				if value.(*HikClient).hikConfig.Ip == v.IPv4Address {
// 					if value.(*HikClient).isConnect && !value.(*HikClient).longConnect {
// 						h.StartLongHikGateway(key.(string))
// 					}
// 				}
// 			}
// 			return true
// 		})
// 	}
// }

// func (h *HikGateway) GetUDPClient() {
// 	//准备广播地址
// 	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:37020")
// 	if err != nil {
// 		// panic(err)
// 		return
// 	}

// 	//准备监听地址
// 	listenAddr, err := net.ResolveUDPAddr("udp4", ":37020")
// 	if err != nil {
// 		return
// 	}

// 	//创建连接
// 	conn, err := net.ListenUDP("udp4", listenAddr)
// 	if err != nil {
// 		return

// 	}
// 	defer func(conn *net.UDPConn) {
// 		err := conn.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(conn)

// 	//向广播地址发送探测数据
// 	uuidString := strings.ToUpper(uuid.NewString())
// 	req := Probe{
// 		Uuid:  uuidString,
// 		Types: "inquiry",
// 	}
// 	sendBytes, err := xml.Marshal(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	_, err = conn.WriteToUDP(sendBytes, addr)
// 	if err != nil {
// 		return
// 	}

// 	//接收回复数据
// 	deviceList := &DeviceList{}
// 	// 判断文件是否存在
// 	// _, err = os.Stat(OUTPUT)
// 	// if os.IsNotExist(err) { // 不存在则创建文件
// 	// 	err := os.WriteFile(OUTPUT, nil, 0644)
// 	// 	if err != nil {
// 	// 		return
// 	// 	}
// 	// }
// 	for {
// 		data := make([]byte, 2048)
// 		// 设置读取超时时间,否则会持续阻塞，一般情况下2-3秒就接收完了，如果时间太短，可按需分配
// 		err := conn.SetReadDeadline(time.Now().Add(time.Second * 2))
// 		if err != nil {
// 			return
// 		}
// 		n, _, err := conn.ReadFromUDP(data)
// 		if err != nil {
// 			return
// 		}
// 		var ipc Device
// 		err = xml.Unmarshal(data[:n], &ipc)
// 		if err != nil {
// 			fmt.Println("xml转换结构体异常：", err.Error())
// 		}
// 		//打印回复数据
// 		deviceList.Lock()
// 		deviceList.Devices = append(deviceList.Devices, ipc)
// 		deviceList.Unlock()
// 		// 编码为JSON格式
// 		// jsonBytes, err := json.MarshalIndent(deviceList.Devices, "", "    ")
// 		h.devices = deviceList.Devices
// 		// if err != nil {
// 		// 	fmt.Println("编码 JSON 时出错:", err)
// 		// 	return
// 		// }
// 		// 写入文件
// 		// err = os.WriteFile(OUTPUT, jsonBytes, 0644)
// 		// if err != nil {
// 		// 	fmt.Println("写入文件错误:", err)
// 		// 	return
// 		// }
// 	}
// }

// 获取udp广播列表
// func (h *HikGateway) GetUDPDeviceList() (resp []Device) {
// 	return h.devices
// }

// 注册一个Hik服务
func (h *HikGateway) RegisterHikGateway(hikConfig HikConfig) (err error) {

	// 判断是否已经注册
	if _, ok := h.hikMap.Load(hikConfig.Ip); ok {
		return errors.New("已经注册")
	}

	hikClient, err := NewOnlyHikClinet(hikConfig)
	if err != nil {
		return
	}

	h.hikMap.Store(hikConfig.Ip, hikClient)

	// 接收广播
	// hikClient.ctx, hikClient.cancel = context.WithCancel(context.Background())
	// go h.handelBoardcast(hikClient)

	// h.longMap.Store(hikConfig.Ip, cancel)

	// // 注册任务服务
	// task := gocron.NewScheduler()
	// task.Every(20).Seconds().Do(h.UserCheck, hikConfig)
	// <-task.Start()
	// h.taskMap.Store(hikConfig, task)
	return
}

// 删除一个Hik服务
func (h *HikGateway) DeleteHikGateway(key string) {

	// 删除服务
	h.cancelHikGateway(key)
}

func (h *HikGateway) cancelHikGateway(key string) {
	if _, ok := h.hikMap.Load(key); ok {
		// 关闭长连接
		h.hikMap.Delete(key)
	}
}

func (h *HikGateway) GetLongAcsEventInfo(key string) (result []EventInfo) {
	if v, ok := h.hikMap.Load(key); ok {
		client := v.(*HikClient)
		if client != nil {
			result = client.doorEventLogs
		}
	}
	return result
}

// 取消后一个hik的长连接
func (h *HikGateway) CancelLongHikGateway(key string) {
	if v, ok := h.hikMap.Load(key); ok {
		// 关闭长连接
		if v.(*HikClient).cancel != nil {
			v.(*HikClient).cancel()

		}
	}
}

// 开启某一个hik的长连接
func (h *HikGateway) StartLongHikGateway(key string) {
	if v, ok := h.hikMap.Load(key); ok {
		// 关闭长连接
		if v.(*HikClient).cancel != nil {
			v.(*HikClient).cancel()
		}
		v.(*HikClient).ctx, v.(*HikClient).cancel = context.WithCancel(context.Background())
		go h.handelBoardcast(v.(*HikClient))
	}
}

// 关闭全部的hik的长连接
func (h *HikGateway) CancelAllLongHikGateway() {
	h.hikMap.Range(func(key, value any) bool {
		// 关闭长连接
		if value.(*HikClient).cancel != nil {
			value.(*HikClient).cancel()
		}
		return true
	})
}

// 开启全部的hik的长连接
func (h *HikGateway) StartAllLongHikGateway() {
	h.hikMap.Range(func(key, value any) bool {
		// 关闭长连接
		if value.(*HikClient).cancel != nil {
			value.(*HikClient).cancel()
		}
		value.(*HikClient).ctx, value.(*HikClient).cancel = context.WithCancel(context.Background())
		go h.handelBoardcast(value.(*HikClient))
		return true
	})
}

// 获取某个服务的连接状态
func (h *HikGateway) GetHikGatewayStatus(key string) (resp HikInfo) {
	if v, ok := h.hikMap.Load(key); ok {
		client := v.(*HikClient)
		if client != nil {
			resp.Ip = v.(*HikClient).hikConfig.Ip
			resp.Port = v.(*HikClient).hikConfig.Port
			resp.Username = v.(*HikClient).hikConfig.Username
			resp.Password = v.(*HikClient).hikConfig.Password
			resp.IsConnect = v.(*HikClient).isConnect
			resp.LongConnect = v.(*HikClient).longConnect
		}
	}
	return
}

// 更新一个Hik服务
func (h *HikGateway) UpdateHikGateway(hikConfig HikConfig) (err error) {
	// 获取原来的服务
	hikConfig.Ip = strings.Replace(hikConfig.Ip, "http://", "", -1)
	hikConfig.Ip = strings.Replace(hikConfig.Ip, "https://", "", -1)
	if _, ok := h.hikMap.Load(hikConfig.Ip); !ok {
		return errors.New("没有该服务")
	} else {
		hikClient := &HikClient{}
		if hikClient, err = newHikClinet(hikConfig); err != nil {
			return err
		}
		// 结束之前的长连接
		// h.DeleteHikGateway(v.(*HikClient).hikConfig.Ip)
		h.hikMap.Store(hikConfig.Ip, hikClient)
		// hikClient.ctx, hikClient.cancel = context.WithCancel(context.Background())
		// go h.handelBoardcast(hikClient)
		// h.longMap.Store(hikConfig.Ip, cancel)
	}

	return
}

// 转发服务
func (h *HikGateway) transferService(key string) (client *HikClient, err error) {
	if v, ok := h.hikMap.Load(key); !ok {
		return nil, errors.New("没有该服务")
	} else {
		return v.(*HikClient), nil
	}
}

// 获取设备连接状态信息
func (h *HikGateway) GetDeviceStatus(key string) (resp []GateWayConnect) {
	h.hikMap.Range(func(key, value any) bool {
		data := GateWayConnect{}
		data.HikConfig = value.(*HikClient).hikConfig
		data.IsConnect = value.(*HikClient).isConnect
		resp = append(resp, data)
		return true
	})
	return
}

// watch全部的广播服务
func (h *HikGateway) handelBoardcast(client *HikClient) {

	if !client.isConnect {
		return
	}
	// file, _ := os.OpenFile(fmt.Sprintf("./%s.txt", client.hikConfig.Ip), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	// defer file.Close()

	client.newHikClientLongConnect(func(buf Result) bool {
		// file.WriteString("接收来自网关的数据\n")
		select {
		case <-client.ctx.Done():
			{
				// file.WriteString("接收来自网关的数据done\n")
				return false
			}
		default:
			{
				// file.WriteString("进入处理网关的数据default\n")
				data := Msg{}
				data.Ip = client.hikConfig.Ip
				data.Type = SUCCESS
				data.Data = buf
				h.broadClient.Range(func(key, value any) bool {
					// file.WriteString("进入处理网关的数据default   发送到通道开始\n")
					key.(chan Msg) <- data
					// file.WriteString(fmt.Sprintf("发送数据向网关注册的服务 发送到通道结束\n=======%s\n%s\n=====\n", buf.Data, buf.Type))
					return true
				})
				// os.Exit(0)
				// file.WriteString("\n开始发送结束")
			}
		}
		return true
	})
	client.isConnect = false

}

// cancel广播服务
func (h *HikGateway) cancelBoardcast(key string) {
	if v, ok := h.hikMap.Load(key); ok {
		v.(*HikClient).cancel()
	}
}

// 注册广播客户端
func (h *HikGateway) RegisterBroadClient(channel chan Msg) {
	h.broadClient.Store(channel, nil)
}

// 删除广播客户端
func (h *HikGateway) DeleteBroadClient(channel chan Msg) {
	h.broadClient.Delete(channel)
}

// 获取网关下的所有设备
func (h *HikGateway) GetGatewayDeviceList() (resp []HikInfo) {
	h.hikMap.Range(func(key, value any) bool {
		resp = append(resp, HikInfo{
			Ip:          value.(*HikClient).hikConfig.Ip,
			Port:        value.(*HikClient).hikConfig.Port,
			Username:    value.(*HikClient).hikConfig.Username,
			Password:    value.(*HikClient).hikConfig.Password,
			IsConnect:   value.(*HikClient).isConnect,
			LongConnect: value.(*HikClient).longConnect,
		})
		return true
	})
	return
}

// 获取缓存门锁状态
func (h *HikGateway) GetDoorStatus(ip string) (result interface{}) {

	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			fmt.Println(value.(*HikClient).doorStatus)
			result = value.(*HikClient).doorStatus
			return false
		}
		return true
	})
	return
}

func (h *HikGateway) GetDoorThreeInfo(ip string) (doorAcsEventLogs *GetAcsEventRes, doorAcsEventTotalNum *GetAcsEventTotalNumRes, doorPersonInfoCount *GetPersonInfoCountRes) {

	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			doorAcsEventLogs = value.(*HikClient).doorAcsEventLogs
			doorAcsEventTotalNum = value.(*HikClient).doorAcsEventTotalNum
			doorPersonInfoCount = value.(*HikClient).doorPersonInfoCount
			return false
		}
		return true
	})
	return

}

// 更改车牌号
func (h *HikGateway) SetVCLData(ip string, data SetVCLDataReq) (err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			err = value.(*HikClient).SetVCLData(data)
			return false
		}
		return true
	})
	return
}

// 设置车牌号
func (h *HikGateway) VCLGetCond(ip string, data SetVCLDataReq) (err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			err = value.(*HikClient).VCLGetCond(data)
			return false
		}
		return true
	})
	return
}

func (h *HikGateway) TCG225EVCLGetCond(ip string, data TCG225EVCLGetCondReq) (err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			err = value.(*HikClient).TCG225EVCLGetCond(data)
			return false
		}
		return true
	})
	return
}

// 删除车牌号
func (h *HikGateway) VCLDelCond(ip string, data VCLDelCondReq) (err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			err = value.(*HikClient).VCLDelCond(data)
			return false
		}
		return true
	})
	return
}

// 删除车牌号
func (h *HikGateway) TCG225EVCLDelCond(ip string, data TCG225EVCLDelCondReq) (err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			err = value.(*HikClient).TCG225EVCLDelCond(data)
			return false
		}
		return true
	})
	return
}

// 获取车牌号列表
func (h *HikGateway) VCLGetList(ip string, data VCLGetListReq) (result *VCLGetListRes, err error) {
	h.hikMap.Range(func(key, value any) bool {
		if key.(string) == ip {
			result, err = value.(*HikClient).VCLGetList(data)
			return false
		}
		return true
	})
	return
}
