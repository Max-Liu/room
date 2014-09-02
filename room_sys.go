package room

type RoomSysInfo struct {
	Pid       int
	Addr      string
	StartTime int64
}

func NewRoom(pid int, addr string, startTime int64) *RoomSysInfo {
	return &RoomSysInfo{
		pid,
		addr,
		startTime,
	}
}
