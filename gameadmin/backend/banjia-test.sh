#!/bin/bash
rsync -avzP --exclude-from=banjia-exclude.txt bin *.sh doudou-test:/data/game/
