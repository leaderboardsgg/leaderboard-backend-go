#! /bin/bash

psql -U postgres -c "CREATE DATABASE leaderboard"
psql -U postgres -d leaderboard -f schema.sql
psql -U postgres -d leaderboard -f testdata.sql
