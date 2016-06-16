###########################################################
#                                                         #
# Name   : stat_zzmm.sh                                   #
#                                                         #
# Usage  :                                                #
#      This script is used to start the Wanbu_Stat_ZZMM   #
#      application.                                       #
#                                                         #
# Author : HuLiWei                                        #
#                                                         #
# History: 2016-03-16 created                             #
#                                                         #
###########################################################

#配置执行文件所在的路径
#工程部署时需改变
PROGRAM_HOME="/home/ylx/piapia/Wanbu_Stat_ZZMM/bin"

#配置APP执行文件的名称
APP="Wanbu_Stat_ZZMM"
LOC=$(cd "$(dirname "$0")"; pwd)
cd $PROGRAM_HOME
start()
{
	echo ""	
	echo "Begin to run $APP ....."
	echo "配置项PROGRAM_HOME为：$PROGRAM_HOME"
	echo "当前模块所在目录：$LOC"
	if [ "$PROGRAM_HOME"x = "$LOC"x ]
	then
		echo "目录配置正确"
		list > /dev/null 2>&1
		if [ $? -ne 0 ] 
		then
			echo "Start failed. The program had been run !"
		else		
			./$APP > /dev/null 2>&1 &
			sleep 1
			ps -ef| grep $APP | grep -v grep >/dev/null 2>&1
			if test $? -ne 0
			then
				echo Start $APP unsuccessfully!
			else
				echo Start $APP successfully!
			fi
 	 	fi
	 	echo""
	else
		echo "目录配置错误，请重新配置"
	fi
}

stop1()
{
	ps -ef| grep $1 | grep -v grep >/dev/null 2>&1
	if test $? -eq 1
	then
		echo Process $1 is not alive !
        else
		proID=`ps -ef| grep $1 |grep -v grep| awk '{ print $2 }'`
		kill -9 $proID 2>/dev/null
		ps -ef| grep $1 | grep -v grep >/dev/null 2>&1
		if test $? -eq 0
		then
			kill -9 $proID 2>/dev/null
		fi
		echo Stop $1 successfully!
       fi
}

stop()
{

	stop1 $APP

	echo ""
}

about()
{
	if [ $APP ]
	then 
		if [ -x $APP ]
		then
			./$APP -about
		else
			echo "There are no file [$APP] here"
		fi
	else 
		./$APP -about
	fi
	echo ""
}

usage()
{
	echo "Usage:" 
	echo "  查看版本号	$0 about"
	echo "  查看启动情况	$0 list"
	echo "  启动 $0 start"
	echo "  停止 $0 stop"
	echo ""
}

list1()
{
	echo "  @$1[$2]:"
	ps -ef| grep $2 | grep -v grep >/dev/null 2>&1
	if test $? -eq 1
	then
		echo "	The process is not alive"
		return 0
	else
		ps -ef | grep $2 |grep -v grep| awk '{print "\t"$1" "$(NF-1)" "$NF}' | sort
		return 1
	fi
}

list()
{
	echo "  >>> The active process(es) as the following:"

	a=0
	list1  $APP模块  $APP
	a=`expr $a + $?`

	echo ""
	echo $a
	return $a
}


# See how we were called.
case "$1" in
	start)		start				2>/dev/null;;
	stop)		stop				2>/dev/null;;
	about)		about 				2>/dev/null;;
	list)		list				2>/dev/null;;
	*)		usage				2>/dev/null;;
esac
exit 0
