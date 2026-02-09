// Copyright (c) Asyraf Norafandi
// SPDX-License-Identifier: MPL-2.0

package client

type Chain struct {
	Slug     string    `json:"slug"`
	Networks []Network `json:"networks"`
}

type Network struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
