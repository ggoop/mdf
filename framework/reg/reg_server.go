package reg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ggoop/mdf/utils"
	"github.com/golang/glog"
	"io"
	"os"
	"path"
)

type RegServer interface {
	Online(item RegObject)
	Offline(item RegObject)
}

var regServerStore *regStoreSv

type regStoreSv struct {
	data   map[string]*RegObject
	dbFile string
}

func GetRegServer() RegServer {
	return regServerStore
}

func StartServer() {
	regServerStore = &regStoreSv{data: make(map[string]*RegObject)}
	regServerStore.load()
	regServerStore.register()
}

func (s *regStoreSv) register() {
	address := utils.DefaultConfig.App.Address
	if address == "" {
		address = fmt.Sprintf("http://127.0.0.1:%s", utils.DefaultConfig.App.Port)
	}
	s.Online(RegObject{
		Code:    utils.DefaultConfig.App.Code,
		Name:    utils.DefaultConfig.App.Name,
		Address: address,
		Configs: utils.DefaultConfig,
	})
}
func (s *regStoreSv) Online(item RegObject) {
	item.Time = utils.TimeNowPtr()
	s.data[item.Key()] = &item
	setRegObjectCache(item.Code, &item)

	s.save()
}
func (s *regStoreSv) Offline(item RegObject) {
	s.data[item.Key()] = nil
	setRegObjectCache(item.Code, nil)

	s.save()
}

func (s *regStoreSv) Get(item RegObject) *RegObject {
	if old, ok := s.data[item.Key()]; ok && old != nil {
		return old
	} else {
		return nil
	}
}

func (s *regStoreSv) GetAll() []*RegObject {
	items := make([]*RegObject, 0)
	for _, item := range s.data {
		items = append(items, item)
	}
	return items
}
func (s *regStoreSv) load() {
	if s.dbFile == "" {
		s.dbFile = utils.JoinCurrentPath(path.Join(utils.DefaultConfig.App.Storage, "system", "regs.db"))
	}
	if !utils.PathExists(s.dbFile) {
		return
	}
	fi, err := os.Open(s.dbFile)
	if err != nil {
		glog.Error(err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		item := RegObject{}
		err = json.Unmarshal(a, &item)
		if err != nil {
			glog.Error(err)
			return
		}
		s.Online(item)
	}
}
func (s *regStoreSv) save() {
	items := s.GetAll()
	f, err := os.Create(s.dbFile)
	if err != nil {
		glog.Error(err)
		f.Close()
		return
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		b, err := json.Marshal(item)
		if err != nil {
			glog.Error(err)
			return
		}
		fmt.Fprintln(f, string(b))
		if err != nil {
			glog.Error(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		glog.Error(err)
		return
	}
}
