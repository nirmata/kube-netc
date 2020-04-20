package tracker

func (t *Tracker) GetConnectionData() map[ConnectionID]ExportData {
	m := make(map[ConnectionID]ExportData)
	for k, v := range t.dataHistory {
		m[k] = ExportData{
			BytesSent:          v.bytesSent,
			BytesRecv:          v.bytesRecv,
			BytesSentPerSecond: v.bytesSentPerSecond,
			BytesRecvPerSecond: v.bytesRecvPerSecond,
			LastUpdated:        v.lastUpdated,
		}
	}
	return m
}

func (t *Tracker) GetNumConnections() uint16 {
	return t.numConnections
}

// Recv

func (t *Tracker) GetBytesRecvPerSecond() uint64 {
	return t.bytesRecvPerSecond
}

func (t *Tracker) GetBytesRecv() uint64 {
	return t.bytesRecv
}

func (t *Tracker) GetTotalBytesRecv() uint64 {
	return t.totalRecv
}

// Sent

func (t *Tracker) GetBytesSentPerSecond() uint64 {
	return t.bytesSentPerSecond
}

func (t *Tracker) GetBytesSent() uint64 {
	return t.bytesSent
}

func (t *Tracker) GetTotalBytesSent() uint64 {
	return t.totalSent
}
