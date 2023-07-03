#!/bin/sh
cd "`echo $0 | rev | cut -c8- | rev`"
chmod +x ./rbq_anonymous_bot_macOS64
# 如果需要添加参数，请在此行后面添加：
# If you need to add parameters, please add after this line:
./rbq_anonymous_bot_macOS64
#
read -p "按 Command+Q 退出。"
exit
