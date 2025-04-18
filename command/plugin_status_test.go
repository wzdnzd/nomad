// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/nomad/ci"
	"github.com/hashicorp/nomad/nomad/state"
	"github.com/posener/complete"
	"github.com/shoenig/test/must"
)

func TestPluginStatusCommand_Implements(t *testing.T) {
	ci.Parallel(t)
	var _ cli.Command = &PluginStatusCommand{}
}

func TestPluginStatusCommand_Fails(t *testing.T) {
	ci.Parallel(t)
	ui := cli.NewMockUi()
	cmd := &PluginStatusCommand{Meta: Meta{Ui: ui}}

	// Fails on misuse
	code := cmd.Run([]string{"some", "bad", "args"})
	must.One(t, code)

	out := ui.ErrorWriter.String()
	must.StrContains(t, out, commandErrorText(cmd))
	ui.ErrorWriter.Reset()

	// Test an unsupported plugin type.
	code = cmd.Run([]string{"-type=not-a-plugin"})
	must.One(t, code)

	out = ui.ErrorWriter.String()
	must.StrContains(t, out, "Unsupported plugin type: not-a-plugin")
	ui.ErrorWriter.Reset()
}

func TestPluginStatusCommand_AutocompleteArgs(t *testing.T) {
	ci.Parallel(t)

	srv, _, url := testServer(t, true, nil)
	defer srv.Shutdown()

	ui := cli.NewMockUi()
	cmd := &PluginStatusCommand{Meta: Meta{Ui: ui, flagAddress: url}}

	// Create a plugin
	id := "long-plugin-id"
	s := srv.Agent.Server().State()
	cleanup := state.CreateTestCSIPlugin(s, id)
	defer cleanup()
	ws := memdb.NewWatchSet()
	plug, err := s.CSIPluginByID(ws, id)
	must.NoError(t, err)

	prefix := plug.ID[:len(plug.ID)-5]
	args := complete.Args{Last: prefix}
	predictor := cmd.AutocompleteArgs()

	res := predictor.Predict(args)
	must.Len(t, 1, res)
	must.Eq(t, plug.ID, res[0])
}
