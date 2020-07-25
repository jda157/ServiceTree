package service_tree

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myprojects/26_concurrency/internal/app/service"
)

type serviceTreeNode struct {
	service *service.Service
	parent  *serviceTreeNode
}

type ServiceTree struct {
	tree map[string]*serviceTreeNode
}

func New() *ServiceTree {
	return &ServiceTree{
		make(map[string]*serviceTreeNode),
	}
}

func (st *ServiceTree) AddService(name string, parent string, handler service.ServiceHandlerT) error {
	if _, ok := st.tree[name]; ok {
		return fmt.Errorf("service %s already exist", name)
	}

	if parent == "" {
		st.tree[name] = &serviceTreeNode{
			service: service.New(name, handler),
			parent:  nil,
		}
		return nil
	}
	if _, ok := st.tree[parent]; !ok {
		return fmt.Errorf("parent %s doesn't  exist", parent)
	}
	st.tree[name] = &serviceTreeNode{
		service: service.New(name, handler),
		parent:  st.tree[parent],
	}
	return nil
}

type parseConf struct {
	Name   string
	Parent string
}

func (st *ServiceTree) InitFromConfig(fileName string) error {
	jsonFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var confs []parseConf
	if err := json.Unmarshal(jsonFile, &confs); err != nil {
		return err
	}

	for _, c := range confs {
		if err := st.AddService(c.Name, c.Parent, service.NewHandler()); err != nil {
			return err
		}
	}
	return nil
}

func (stn *serviceTreeNode) solveDependency(ctx context.Context) (service.OutChanT, error) {
	if stn.parent != nil {
		out, err := stn.parent.solveDependency(ctx)
		if err != nil {
			return nil, err
		}
		return stn.service.CallHandler(ctx, service.InChanT{true, out})
	}

	return stn.service.CallHandler(ctx, service.InChanT{false, nil})
}

func (st *ServiceTree) CallService(ctx context.Context, name string) (int, error) {
	if _, ok := st.tree[name]; !ok {
		return 0, fmt.Errorf("service '%s'  doesn't exist", name)
	}

	out, err := st.tree[name].solveDependency(ctx)
	if err != nil {
		return -1, err
	}
	select {
	case <-ctx.Done():
		return -1, service.ErrTimeout
	case val := <-out:
		return val, nil
	}
}

