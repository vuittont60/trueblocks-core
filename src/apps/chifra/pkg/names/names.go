package names

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type Parts int

// Parts is a bitfield that defines what parts of a name to return and other options
const (
	None      Parts = 0x0
	Regular   Parts = 0x1
	Custom    Parts = 0x2
	Prefund   Parts = 0x4
	Testing   Parts = 0x8
	MatchCase Parts = 0x10
	Expanded  Parts = 0x20
)

type SortBy int

// SortBy is a bitfield that defines how to sort the names
const (
	SortByAddress SortBy = iota
	SortByName
	// SortBySymbol
	// SortBySource
	// SortByDecimals
	SortByTags
	// SortByPetname
)

// LoadNamesArray loads the names from the cache and returns an array of names
func LoadNamesArray(chain string, parts Parts, sortBy SortBy, terms []string) ([]types.SimpleName, error) {
	var names []types.SimpleName
	if namesA, err := LoadNamesMap(chain, parts, terms); err != nil {
		return nil, err
	} else {
		for _, name := range namesA {
			isTesting := parts&Testing != 0
			isIndiv := strings.Contains(name.Tags, "Individual")
			if name.Address.Hex() == "0x69e271483c38ed4902a55c3ea8aab9e7cc8617e5" {
				isIndiv = false
				name.Name = "Name 0x69e27148"
			}
			if !isTesting || !isIndiv {
				names = append(names, name)
			}
		}
	}

	sort.Slice(names, func(i, j int) bool {
		switch sortBy {
		case SortByName:
			return names[i].Name < names[j].Name
		case SortByTags:
			return names[i].Tags < names[j].Tags
		case SortByAddress:
			fallthrough
		default:
			return names[i].Address.Hex() < names[j].Address.Hex()
		}
	})

	isTesting := parts&Testing != 0
	isTags := sortBy == SortByTags
	if isTesting && !isTags {
		names = names[:utils.Min(200, len(names))]
	}

	return names, nil
}

// LoadNamesMap loads the names from the cache and returns a map of names
func LoadNamesMap(chain string, parts Parts, terms []string) (map[types.Address]types.SimpleName, error) {
	ret := map[types.Address]Name{}

	// Load the prefund names first...
	if parts&Prefund != 0 {
		loadPrefundMap(chain, terms, parts, &ret)
	}

	// binPath := config.GetPathToCache(chain) + "names/names.bin"
	namesPath := filepath.Join(config.GetPathToChainConfig(chain), "names.tab")
	customPath := filepath.Join(config.GetPathToChainConfig(chain), "names_custom.tab")

	// Load the names from the binary file (note that these may overwrite the prefund names)
	if parts&Regular != 0 {
		// enabled := false // os.Getenv("FAST") == "true" // TODO: this isn't right
		// if enabled && file.FileExists(binPath) {
		// 	file, _ := os.OpenFile(binPath, os.O_RDONLY, 0)
		// 	defer file.Close()

		// 	header := NameOnDiscHeader{}
		// 	err := binary.Read(file, binary.LittleEndian, &header)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	for i := uint64(0); i < header.Count; i++ {
		// 		v := NameOnDisc{}
		// 		binary.Read(file, binary.LittleEndian, &v)
		// 		n := Name{
		// 			Tags:       asString("tags", v.Tags[:]),
		// 			Address:    asString("address", v.Address[:]),
		// 			Name:       asString("name", v.Name[:]),
		// 			Symbol:     asString("symbol", v.Symbol[:]),
		// 			Decimals:   fmt.Sprintf("%d", v.Decimals),
		// 			Source:     asString("source", v.Source[:]),
		// 			Petname:    asString("petname", v.Petname[:]),
		// 			IsCustom:   v.Flags&IsCustom != 0,
		// 			IsPrefund:  v.Flags&IsPrefund != 0,
		// 			IsContract: v.Flags&IsContract != 0,
		// 			IsErc20:    v.Flags&IsErc20 != 0,
		// 			IsErc721:   v.Flags&IsErc721 != 0,
		// 			Deleted:    v.Flags&IsDeleted != 0,
		// 		}
		// 		if !n.IsCustom {
		// 			if doSearch(n, terms, parts) {
		// 				ret[types.HexToAddress(n.Address)] = n
		// 			}
		// 		}
		// 	}
		// } else {
		nameMapFromFile(chain, namesPath, terms, parts, &ret)
		// }
	}

	// Load the custom names (note that these may overwrite the prefund and regular names)
	if parts&Custom != 0 {
		loadCustomMap(chain, customPath, terms, parts, &ret)
	}

	ret2 := map[types.Address]types.SimpleName{}
	for k, v := range ret {
		ret2[k] = v.ToSimpleName()
	}

	return ret2, nil
}

// Name is a record in the names database
type Name struct {
	Tags       string `json:"tags"`
	Address    string `json:"address"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Source     string `json:"source"`
	Decimals   string `json:"decimals"`
	Petname    string `json:"petname"`
	Deleted    bool   `json:"deleted"`
	IsCustom   bool   `json:"isCustom"`
	IsPrefund  bool   `json:"isPrefund"`
	IsContract bool   `json:"isContract"`
	IsErc20    bool   `json:"isErc20"`
	IsErc721   bool   `json:"isErc721"`
}

func (n Name) String() string {
	ret, _ := json.MarshalIndent(n, "", "  ")
	return string(ret)
}

var requiredColumns = []string{
	"tags",
	"address",
	"name",
	"symbol",
	"source",
	"petname",
}

type NameReader struct {
	file      *os.File
	header    map[string]int
	csvReader csv.Reader
}

func (gr *NameReader) Read() (Name, error) {
	record, err := gr.csvReader.Read()
	if err == io.EOF {
		gr.file.Close()
	}
	if err != nil {
		return Name{}, err
	}

	return Name{
		Tags:       record[gr.header["tags"]],
		Address:    strings.ToLower(record[gr.header["address"]]),
		Name:       record[gr.header["name"]],
		Decimals:   record[gr.header["decimals"]],
		Symbol:     record[gr.header["symbol"]],
		Source:     record[gr.header["source"]],
		Petname:    record[gr.header["petname"]],
		IsCustom:   record[gr.header["iscustom"]] == "true",
		IsPrefund:  record[gr.header["isprefund"]] == "true",
		IsContract: record[gr.header["iscontract"]] == "true",
		IsErc20:    record[gr.header["iserc20"]] == "true",
		IsErc721:   record[gr.header["iserc721"]] == "true",
		Deleted:    record[gr.header["deleted"]] == "true",
	}, nil
}

func NewNameReader(path string) (NameReader, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return NameReader{}, err
	}

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	if strings.HasSuffix(path, ".csv") {
		reader.Comma = ','
	}

	headerRow, err := reader.Read()
	if err != nil {
		return NameReader{}, err
	}
	header := map[string]int{}
	for index, columnName := range headerRow {
		header[columnName] = index
	}

	for _, required := range requiredColumns {
		_, ok := header[required]
		if !ok {
			err = fmt.Errorf(`required column "%s" missing in file %s`, required, path)
			return NameReader{}, err
		}
	}

	gr := NameReader{
		file:      file,
		header:    header,
		csvReader: *reader,
	}

	return gr, nil
}

func (n *Name) ToSimpleName() types.SimpleName {
	return types.SimpleName{
		Tags:       n.Tags,
		Address:    types.HexToAddress(n.Address),
		Name:       n.Name,
		Source:     n.Source,
		Petname:    n.Petname,
		Symbol:     n.Symbol,
		Decimals:   globals.ToUint64(n.Decimals),
		Deleted:    n.Deleted,
		IsCustom:   n.IsCustom,
		IsPrefund:  n.IsPrefund,
		IsContract: n.IsContract,
		IsErc20:    n.IsErc20,
		IsErc721:   n.IsErc721,
	}
}

/*
// NameOnDiscHeader is the header of the names database when stored in the binary backing file
type NameOnDiscHeader struct {
	Magic   uint64    // 8 bytes
	Version uint64    // + 8 bytes = 16 bytes
	Count   uint64    // + 8 bytes = 24 bytes
	Padding [428]byte // 452 - 24 = 428 bytes
}

// NameOnDisc is a record in the names database when stored in the binary backing file
type NameOnDisc struct {
	Tags     [30 + 1]byte  `json:"-"` // 31 bytes
	Address  [42 + 1]byte  `json:"-"` // + 43 bytes = 74 bytes
	Name     [120 + 1]byte `json:"-"` // + 121 bytes = 195 bytes
	Symbol   [30 + 1]byte  `json:"-"` // + 31 bytes = 226 bytes
	Source   [180 + 1]byte `json:"-"` // + 181 bytes = 407 bytes
	Petname  [40 + 1]byte  `json:"-"` // + 41 bytes = 448 bytes
	Decimals uint16        `json:"-"` // + 2 bytes = 450 bytes
	Flags    uint16        `json:"-"` // + 2 bytes = 452 bytes
}

// Bitflags for the Flags field
const (
	IsCustom   uint16 = 0x1
	IsPrefund  uint16 = 0x2
	IsContract uint16 = 0x4
	IsErc20    uint16 = 0x8
	IsErc721   uint16 = 0x10
	IsDeleted  uint16 = 0x20
)

// asString converts the byte array (not zero-terminated) to a string
func asString(which string, b []byte) string {
	ret := ""
	for _, rVal := range string(b) {
		if rVal == 0 {
			return ret
		}
		ret += string(rVal)
	}
	return ret
}
*/
