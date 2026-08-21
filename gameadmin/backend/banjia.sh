#!/bin/bash
# rsync -n -avzP --exclude-from=banjia-exclude.txt bin *.sh doudou-prod:/data/game/
#rsync  -avzP --exclude-from=banjia-exclude.txt bin *.sh dou-ph-prod:/data/game/
rsync  -avzP --exclude-from=banjia-exclude.txt bin *.sh dou-idr-prod:/data/game/
