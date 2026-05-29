#!/usr/bin/env bash
# Copyright Adfinis 2025, 2026
# SPDX-License-Identifier: GPL-3.0-or-later


mkdir -p ssh-keys
ssh-keygen -t ed25519 -f ssh-keys/id_ed25519 -N "" -C "bastion-test-key"
