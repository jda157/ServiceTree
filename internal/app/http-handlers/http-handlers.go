package http_handlers

import (
	"context"
	"encoding/json"
	"myprojects/26_concurrency/internal/app/service"
	service_tree "myprojects/26_concurrency/internal/app/service-tree"
	"net/http"
	"time"
)

type chanStructT struct {
	err error
	val int
}

type responseStruct struct {
	ServiceName string
	Status      string
	Data        int
}

func HttpHandler(st *service_tree.ServiceTree) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var services []string
		if err := json.NewDecoder(r.Body).Decode(&services); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		channelsResp := make([]chan chanStructT, len(services))

		for i, service := range services {
			channelsResp[i] = make(chan chanStructT)
			go func(name string, num int) {
				val, err := st.CallService(ctx, name)
				if err != nil {
					channelsResp[num] <- chanStructT{
						err: err,
					}
					return
				}
				channelsResp[num] <- chanStructT{
					err: nil,
					val: val,
				}
			}(service, i)
		}

		responses := make([]responseStruct, len(services))
		for num, channel := range channelsResp {
			resp := <-channel
			responses[num].ServiceName = services[num]
			if resp.err != nil {
				if resp.err == service.ErrTimeout {
					responses[num].Status = "Timeout"
				} else {
					responses[num].Status = "Error"
				}
			} else {
				responses[num].Status = "Ok"
				responses[num].Data = resp.val
			}
		}
		data, err := json.MarshalIndent(&responses, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
