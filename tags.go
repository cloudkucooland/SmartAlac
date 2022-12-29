package sa

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abema/go-mp4"
)

func LoadALACTags(f string) (map[string]*mp4.Data, error) {
	info := make(map[string]*mp4.Data)
	mean := "com.apple.iTunes" // seems to always be this
	name := "unkn"
	buf := new(bytes.Buffer)

	in, err := os.OpenFile(f, os.O_RDONLY, 0755)
	if err != nil {
		return info, err
	}
	defer in.Close()

	_, err = mp4.ReadBoxStructure(in, func(h *mp4.ReadHandle) (interface{}, error) {
		switch h.BoxInfo.Type.String() {
		case "mean":
			if !h.BoxInfo.Context.UnderIlstMeta {
				log.Println("mean not under ilst", h.Path)
				return h.Expand()
			}
			_, err = h.ReadData(buf)
			if err != nil {
				return nil, err
			}
			mean = buf.String()
			buf.Reset()
			return h.Expand()
		case "name":
			if !h.BoxInfo.Context.UnderIlstMeta {
				log.Println("name not under ilst", h.Path)
				return h.Expand()
			}
			_, err = h.ReadData(buf)
			if err != nil {
				return nil, err
			}
			name = buf.String()
			buf.Reset()
			return h.Expand()
		case "----":
			// expand to get mean/name
			return h.Expand()
		case "data":
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}

			p := h.Path[len(h.Path)-2].String()
			if p == "----" {
				p = fmt.Sprintf("%s.%s", mean, name)
			}

			info[p] = box.(*mp4.Data)
			return nil, nil
		default:
			if h.BoxInfo.IsSupportedType() {
				return h.Expand()
			} else {
				log.Println("non-supported type", h.BoxInfo.Type.String())
				return nil, nil
			}
		}
	})
	if err != nil {
		log.Println(err.Error())
		return info, err
	}

	dumpAlac(info)
	return info, nil
}

func dumpAlac(info map[string]*mp4.Data) {
	for k, v := range info {
		switch v.DataType {
		case mp4.DataTypeBinary, mp4.DataTypeSignedIntBigEndian:
			fmt.Printf("%s:\t%d\n", k, binary.BigEndian.Uint32(v.Data))
		case mp4.DataTypeStringUTF8, mp4.DataTypeStringUTF16, mp4.DataTypeStringMac:
			fmt.Printf("%s:\t%s\n", k, strings.TrimSpace(string(v.Data)))
		case mp4.DataTypeFloat32BigEndian:
			// fmt.Printf("%s:\t%d\n", k, binary.BigEndian.Float32(v.Data))
		case mp4.DataTypeFloat64BigEndian:
			// fmt.Printf("%s:\t%d\n", k, binary.BigEndian.Float64(v.Data))
		default:
			fmt.Printf("%s:\t%+v\n", k, v)
		}
	}
}

func SaveLACTags(in string, out string, t map[string]*mp4.Data) error {
	/*
	    open in and out...

		   _, err = mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		   	switch h.BoxInfo.Type {
		   	case mp4.BoxTypeIlst():
		   		ilst, err := w.StartBox(...)
		   		if err != nil {
		   			return nil, err
		   		}

		   		if _, err := mp4.Marshal(w, ilst, h.BoxInfo.Context); err != nil {
		   			return nil, err
		   		}
		   		// rewrite box size
		   		_, err = w.EndBox()
		   		return nil, err
		   	default:
		   		// copy all
		   		return nil, w.CopyBox(r, &h.BoxInfo)
		   	}
		   })
	*/
    return nil
}
