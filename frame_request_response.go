package rsocket

import "io"

type FrameRequestResponse struct {
	*Header
	metadata []byte
	data     []byte
}

func (p *FrameRequestResponse) Size() int {
	size := headerLen
	if p.Header.Flags().Check(FlagMetadata) {
		size += 3 + len(p.metadata)
	}
	if p.data != nil {
		size += len(p.data)
	}
	return size
}

func (p *FrameRequestResponse) WriteTo(w io.Writer) (n int64, err error) {
	var wrote int
	wrote, err = w.Write(p.Header.Bytes())
	n += int64(wrote)
	if err != nil {
		return
	}
	if p.Header.Flags().Check(FlagMetadata) {
		wrote, err = w.Write(encodeU24(len(p.metadata)))
		n += int64(wrote)
		if err != nil {
			return
		}
		wrote, err = w.Write(p.metadata)
		n += int64(wrote)
		if err != nil {
			return
		}
	}
	if p.data == nil {
		return
	}
	wrote, err = w.Write(p.data)
	n += int64(wrote)
	return
}

func (p *FrameRequestResponse) Metadata() []byte {
	return p.metadata
}

func (p *FrameRequestResponse) Data() []byte {
	return p.data
}

func (p *FrameRequestResponse) Parse(h *Header, bs []byte) error {
	p.Header = h
	p.metadata, p.data = sliceMetadataAndData(p.Header, bs, headerLen)
	return nil
}

func mkRequestResponse(sid uint32, metadata []byte, data []byte, f ...Flags) *FrameRequestResponse {
	h := mkHeader(sid, REQUEST_RESPONSE, f...)
	return &FrameRequestResponse{
		Header:   h,
		metadata: metadata,
		data:     data,
	}
}
