#!/usr/bin/env bash
# Copyright (c) Adfinis
# SPDX-License-Identifier: GPL-3.0-or-later


docker compose exec -T bastion \
    /opt/bastion/bin/admin/setup-first-admin-account.sh \
    bastionadmin \
    auto < ssh-keys/id_ed25519.pub
