package tracker

func (t *Tracker) GetConnectionData() map[ConnectionID]ExportData {
	m := make(map[ConnectionID]ExportData)
	for k, v := range t.dataHistory {
		m[k] = ExportData{
			BytesSent: v.bytesSent,
			BytesRecv: v.bytesRecv,
			LastUpdated: v.lastUpdated,
		}
	}
	return m
}

func (t *Tracker) GetNumConnections() uint16 {
	return t.numConnections
}

func (t *Tracker) GetBytesRecvPerSecond() uint64 {
	return t.bytesRecvPerSecond
}

func (t *Tracker) GetBytesRecv() uint64 {
	return t.bytesRecv
}

func (t *Tracker) GetTotalBytesRecv() uint64 {
	return t.totalRecv
}
