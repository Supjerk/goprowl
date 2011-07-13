//
// Copyright (c) 2011, Yanko D Sanchez Bolanos
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the author nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package goprowl

import (
	"http"
	"fmt"
)

const (
	API_SERVER = "https://api.prowlapp.com"
	ADD_PATH   = "/publicapi/add"
	API_URL    = API_SERVER + ADD_PATH
)


type Goprowl struct {
	apikeys []string
}

func (gp *Goprowl) RegisterKey(key string) {

	if len(key) != 40 {

		fmt.Printf("Error, Apikey must be 40 characters long.\n")
		// need to raise an error.
	}

	gp.apikeys = append(gp.apikeys, key)

}

func (gp *Goprowl) DelKey(key string) {
}

func (gp *Goprowl) Push(app string, event string, description string, priority string) {

	ch := make(chan string)

	for _, apikey := range gp.apikeys {

		apikeyList := []string{apikey}
		applicationList := []string{app}
		eventList := []string{event}
		descriptionList := []string{description}
		priorityList := []string{priority}
		vals := http.Values{"apikey": apikeyList,
			"application": applicationList,
			"description": descriptionList,
			"event":       eventList,
			"priority":    priorityList}

		// overkill?
		go func(key string) {
			r, err := http.PostForm(API_URL, vals)

			if err != nil {
				fmt.Printf("%s\n", err)
				ch <- key
			} else {
				if r.StatusCode != 200 {
					ch <- key
				} else {
					ch <- ""
				}
			}

		}(apikey)

	}


	//fmt.Printf("Waiting...\n")
	for i := 0; ; i++ {

		if i == len(gp.apikeys) {
			break
		}

		rc := <-ch
		if rc != "" {
			fmt.Printf("The following key failed: %s\n", rc)
		}

	}

}
