package sa

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
			_, err = h.ReadData(buf)
			if err != nil {
				return nil, err
			}
			mean = buf.String()
			buf.Reset()
			return h.Expand()
		case "name":
			_, err = h.ReadData(buf)
			if err != nil {
				return nil, err
			}
			name = buf.String()
			buf.Reset()
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

			// reset name/mean
			mean = "com.apple.iTunes"
			name = "unkn"
			return nil, nil
		case "----": // remove, just use default
			// expand to get mean/name
			return h.Expand()
		default:
			if h.BoxInfo.IsSupportedType() {
				return h.Expand()
			}
			// log.Println("non-supported type", h.BoxInfo.Type.String())
			return nil, nil
		}
	})
	if err != nil {
		log.Println(err.Error())
		return info, err
	}

	// dumpAlac(info)
	return info, nil
}

func dumpAlac(info map[string]*mp4.Data) {
	for k, v := range info {
		// fmt.Println(v.DataType)
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
			fmt.Printf("%s:\t%+v\n", k, strings.TrimSpace(string(v.Data)))
		}
	}
}

func SaveALACTags(f string, t map[string]*mp4.Data) error {
	in, err := os.OpenFile(f, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	tmp, err := os.CreateTemp("", "smartalac")
	if err != nil {
		log.Fatal(err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	out := mp4.NewWriter(tmp)
	// padsize := 0

	// XXX this will not walk down [moov udta meta ilst ---- data] to get to the data
	_, err = mp4.ReadBoxStructure(in, func(h *mp4.ReadHandle) (interface{}, error) {
		switch h.BoxInfo.Type {
		case mp4.BoxTypeIlst():
			_, err := out.StartBox(&h.BoxInfo)
			if err != nil {
				return nil, err
			}
			ilst := mapToIlst(t, out)
			if _, err := mp4.Marshal(out, ilst, h.BoxInfo.Context); err != nil {
				return nil, err
			}
			// rewrite box size
			_, err = out.EndBox()
			return nil, err
		case mp4.BoxTypeMoov(), mp4.BoxTypeUdta(), mp4.BoxTypeMeta():
			// build the combined box...
			// track size with padsize...
			sub, err := h.Expand()
			return sub, err // for now
		case mp4.BoxTypeFree():
			// XXX figure out how much padding to add
			return nil, out.CopyBox(in, &h.BoxInfo)
		default:
			// copy all
			return nil, out.CopyBox(in, &h.BoxInfo)
		}
	})

	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}

	// mv tmpName to f
	return nil
}

func mapToIlst(t map[string]*mp4.Data, out *mp4.Writer) *mp4.Ilst {
	i := mp4.Ilst{}

	for k, v := range t {
		switch k {
		case "----":
			fmt.Printf("%s:\t%+v\n", k, v)
		default:
			fmt.Printf("%s:\t%+v\n", k, string(v.Data))
		}
	}

	return &i
}

// from https://github.com/Sorrow446/go-mp4tag/blob/main/mp4tag.go
// convert this to some structs to marshall...
func writeCustomMeta(out *mp4.Writer, ctx mp4.Context, val *mp4.Data, field string) error {
	_, err := out.StartBox(&mp4.BoxInfo{Type: mp4.BoxType{'-', '-', '-', '-'}, Context: ctx})
	if err != nil {
		return err
	}
	_, err = out.StartBox(&mp4.BoxInfo{Type: mp4.BoxType{'m', 'e', 'a', 'n'}, Context: ctx})
	if err != nil {
		return err
	}
	_, err = out.Write([]byte{'\x00', '\x00', '\x00', '\x00'})
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, "com.apple.iTunes")
	if err != nil {
		return err
	}
	_, err = out.EndBox()
	if err != nil {
		return err
	}
	_, err = out.StartBox(&mp4.BoxInfo{Type: mp4.BoxType{'n', 'a', 'm', 'e'}, Context: ctx})
	if err != nil {
		return err
	}
	// _, err = out.Write(val.DataLang) // flip to bytes
	_, err = out.Write([]byte{'\x00', '\x00', '\x00', '\x00'})
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, field)
	if err != nil {
		return err
	}
	_, err = out.EndBox()
	if err != nil {
		return err
	}
	/* err = marshalData(out, ctx, val)
	if err != nil {
		return err
	} */
	_, err = out.EndBox()
	return err
}
