// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package model

// Source represents display or publishing preferences of user's own account.
// Returned as an additional entity when verifying and updated credentials, as an attribute of Account.
type Source struct {
	// The default post privacy to be used for new statuses.
	//    public = Public post
	//    unlisted = Unlisted post
	//    private = Followers-only post
	//    direct = Direct post
	Privacy Visibility `json:"privacy"`
	// Whether new statuses should be marked sensitive by default.
	Sensitive bool `json:"sensitive"`
	// The default posting language for new statuses.
	Language string `json:"language"`
	// The default posting content type for new statuses.
	StatusContentType string `json:"status_content_type"`
	// Profile bio.
	Note string `json:"note"`
	// Metadata about the account.
	Fields []Field `json:"fields"`
	// The number of pending follow requests.
	FollowRequestsCount int `json:"follow_requests_count"`
}
