/*
 Copyright (C) 2018 OpenControl Contributors. See LICENSE.md for license.
*/

package resources

import (
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
	"github.com/opencontrol/compliance-masonry/tools/vcs"
)

//go:generate mockery -name Downloader

// Downloader is a generic interface for how to download entries.
type Downloader interface {
	DownloadRepo(common.RemoteSource, string) error
}

// NewVCSDownloader is a constructor for downloading entries using VCS methods.
func NewVCSDownloader() Downloader {
	return vcsEntryDownloader{vcs.Manager{}}
}

type vcsEntryDownloader struct {
	manager vcs.RepoManager
}

// DownloadEntry is a implementation for downloading entries using VCS methods.
func (v vcsEntryDownloader) DownloadRepo(entry common.RemoteSource, destination string) error {
	err := v.manager.Clone(entry.GetURL(), entry.GetRevision(), destination)
	if err != nil {
		return err
	}
	return nil
}
