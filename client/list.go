package client

import (
	"strings"

	"go.etcd.io/etcd/clientv3"
	mvccpb "go.etcd.io/etcd/mvcc/mvccpb"
)

// list a directory
func (clt *EtcdHRCHYClient) List(key string) ([]*Node, error) {
	key, _, err := clt.ensureKey(key)
	if err != nil {
		return nil, err
	}

	// directory start with /
	dir := key + "/"
	if key == "" {
		dir = ""
	}

	// txn := clt.client.Txn(clt.ctx)
	// // make sure the list key is a directory
	// txn.If(
	// 	clientv3.Compare(
	// 		clientv3.Value(key),
	// 		"=",
	// 		clt.dirValue,
	// 	),
	// ).Then(
	// 	clientv3.OpGet(dir, clientv3.WithPrefix()),
	// )

	// txnResp, err := txn.Commit()

	// if !txnResp.Succeeded {
	// 	return nil, ErrorListKey
	// } else {
	// 	if len(txnResp.Responses) > 0 {
	// 		rangeResp := txnResp.Responses[0].GetResponseRange()
	// 		return clt.list(dir, rangeResp.Kvs)
	// 	} else {
	// 		// empty directory
	// 		return []*Node{}, nil
	// 	}
	// }

	txnResp, err := clt.client.Get(clt.ctx, dir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	return clt.list(dir, txnResp.Kvs)
}

// pick key/value under the dir
func (clt *EtcdHRCHYClient) list(dir string, kvs []*mvccpb.KeyValue) ([]*Node, error) {
	nodes := []*Node{}
	// for _, kv := range kvs {
	// 	name := strings.TrimPrefix(string(kv.Key), dir)
	// 	if strings.Contains(name, "/") {
	// 		// secondary directory
	// 		continue
	// 	}
	// 	nodes = append(nodes, clt.createNode(kv))
	// }

	for _, kv := range kvs {
		name := strings.TrimPrefix(string(kv.Key), dir)
		//获取/后的第一个值 bb/cc => bb
		dirName := strings.Split(name, "/")[0]
		len := len(strings.Split(name, "/"))
		isexist := false
		for _, info := range nodes {
			if info.DirName == dirName {
				if info.IsDir {
					info.Count++
					if len > 1 {
						isexist = true
						break
					}
				}
			}
		}
		if !isexist {
			node := &Node{
				DirName:  dirName,
				KeyValue: kv,
				Count:    1,
				IsDir:    len > 1,
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, nil

	// for _, kv := range kvs {
	// 	//拿掉前缀，类似 aa/bb/cc => bb/cc
	// 	name := strings.TrimPrefix(string(kv.Key), dir)
	// 	//获取/后的第一个值 bb/cc => bb
	// 	dirName := strings.Split(name, "/")[0]
	// 	isExist := false
	// 	//遍历nodes
	// 	for _, node := range nodes {
	// 		//判断node是否已经存在，则结束遍历nodes
	// 		if strings.Contains(string(node.Key), dirName) {
	// 			isExist = true
	// 			break
	// 		}
	// 	}
	// 	//不存在的话，新增node
	// 	if !isExist {
	// 		isdir := false
	// 		if strings.Contains(name, "/") {
	// 			isdir = true
	// 		}
	// 		nodes = append(nodes, clt.createNode(kv, dirName, isdir))
	// 	}
	// }

	// return nodes, nil
}
