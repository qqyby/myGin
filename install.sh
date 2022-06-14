#!/bin/bash

#################################
#user options: ./install.sh -2047
#################################
PORT=8081
echo "arg0: port"
echo "for example, $0"
if [[ $1 != "" ]]; then PORT=$1; fi

INSTALL="/opt/bravo/bravo_gin"
BASELICENSE="/opt/bravo/license_files"

if [[ $PORT == "" ]]; then
    SERVICE=bravo_gin
    TARGET_INSTLL=${INSTALL}
else
    SERVICE=bravo_gin
    TARGET_INSTALL=${INSTALL}
fi

INITD=/etc/init.d/${bravo_gin}
LIB=/usr/local/lib

ok_msg(){
    echo -e "${1}${POS}${BLACK}[${GREEN}  OK  ${BLACK}]"
}

failed_msg(){
    echo -e "${1}${POS}${BLACK}[${RED}FAILED${BLACK}]"
}

echo "argv[0]=$0"
if [[ ! -f $0 ]]; then
    echo "directly execute the scripts on shell.";
    work_dir=`pwd`
else
    echo "execute scripts in file: $0";
    work_dir=`dirname $0`; work_dir=`(cd ${work_dir} && pwd)`
fi

# require sudo users
sudo echo "ok";
ret=$?; if [[ 0 -ne ${ret} ]]; then echo "[error]: you must be sudoer"; exit 1; fi

linux_dist_id=`lsb_release -is 2>/dev/null`
ret=$?; if [[ $ret -ne 0 ]];then
    linux_dist_id=`cat /etc/redhat-release 2>/dev/null|awk '{print $1}'`
fi
if [[ 'CentOS' != ${linux_dist_id} ]]; then failed_msg "only support Centos, your os is ${linux_dist_id}"; return 1; fi

linux_dist="$linux_dist_id$linux_dist_main_version"
ok_msg "check os version ${linux_dist_id} ${linux_dist_version} (${linux_dist}) success"


function install(){
    if [[ -f /etc/init.d/ ]]; then
        sudo /etc/init.d/${SERVICE} status >/dev/null 2>&1
        ret=$?; if [[ 0 -eq ${ret} ]]; then
            failed_msg "you must stop the service first: sudo /etc/init.d/${SERVICE} stop";
            return 1;
        fi
    fi

    install_root=$TARGET_INSTALL
    install_bin=$install_root/bravo_gin
    conf_file="$install_root/configs/config.yaml"
    log_file="$install_root/logs"
    sys_log_file="$install_root/objs/sys.log"
    license_file_dir="$BASELICENSE/bravo_transcoder"

    if [[ -d $install_root ]]; then
        version="unknown"
        if [[ -f $install_bin ]]; then
            version=`$install_bin -version`
        fi

        backup_dir=${install_root}.`date "+%Y-%m-%d_%H-%M-%S"`.v-$version
        echo "backup intalled dir, version=$version"
        echo " to=$backup_dir"
        sudo mv $install_root $backup_dir
        ret=$?; if [[ 0 -ne ${ret} ]]; then failed_msg "backup installed dir failed"; return $ret; fi
        ok_msg "backup installed dir success"
    fi

    echo "create install dir"
    sudo mkdir -p $install_root
    ret=$?; if [[ 0 -ne ${ret} ]]; then failed_msg "create install dir failed"; return $ret; fi
    ok_msg "create install dir success"

    ok_msg "prepare compontents"
    (
        ok_msg "prepare init.d script" &&
        sed -i "s|port=.*$|port=${PORT}|g" $work_dir/etc/init.d/bravo_gin &&
        sed -i "s|ROOT=.*$|ROOT=\"${install_root}\"|g" $work_dir/etc/init.d/bravo_gin &&
        sed -i "s|APP=.*$|APP=\"${install_bin}\"|g" $work_dir/etc/init.d/bravo_gin &&
        sed -i "s|CONFIG=.*$|CONFIG=\"${conf_file}\"|g" $work_dir/etc/init.d/bravo_gin &&
        sed -i "s|SYSLOG=.*$|SYSLOG=\"${sys_log_file}\"|g" $work_dir/etc/init.d/bravo_gin &&
        sed -i "s|  run_mode: .*|  run_mode: release|g" $work_dir/configs/config.yaml &&
        sed -i "s|  port: .*|  port: ${PORT}|g" $work_dir/configs/config.yaml &&
        sed -i "s|  license_file_dir: .*|  license_file_dir: ${license_file_dir}|g" $work_dir/configs/config.yaml &&
        sed -i "s|ExecStart=.*$|ExecStart=${INITD} start|g" $work_dir/etc/init.d/bravo_gin.service &&
        sed -i "s|ExecReload=.*$|ExecReload=${INITD} restart|g" $work_dir/etc/init.d/bravo_gin.service &&
        sed -i "s|ExecStop=.*$|ExecStop=${INITD} stop|g" $work_dir/etc/init.d/bravo_gin.service
    )

    echo "copy main file"
    (
        sudo mkdir -p $install_root/logs &&
        sudo cp -rf $work_dir/objs $install_root &&
        sudo cp -f $work_dir/bravo_gin $install_root &&
        sudo cp -rf $work_dir/configs $install_root &&
        sudo cp -rf $work_dir/docs $install_root &&
        sudo mkdir -p $license_file_dir &&
        sudo cp -rf $work_dir/etc $install_root
    )
    ret=$?; if [[ 0 -ne ${ret} ]]; then failed_msg "copy main file failed"; return $ret; fi
    ok_msg "copy main file success"

    ok_msg "install init.d script"
    (
        sudo rm -f ${INITD} &&
        sudo ln -sf $install_root/etc/init.d/bravo_gin ${INITD} &&

        if [[ "CentOS6" == ${linux_dist} ]]; then
            sudo /sbin/chkconfig --add ${INITD}
        else
            sudo rm -rf /lib/systemd/system/bravo_gin.service &&
            sudo cp $work_dir/etc/init.d/bravo_gin.service /lib/systemd/system/ &&
            sudo chmod 754 /lib/systemd/system/bravo_gin.service &&
            sudo systemctl enable bravo_gin.service
        fi
    )
    ret=$?; if [[ 0 -ne ${ret} ]]; then failed_msg "install init.d scripts failed"; return $ret; fi
    ok_msg "install init.d scripts success"
}

install
ret=$?; if [[ $ret -ne 0 ]]; then
    failed_msg " install bravo_gin failed."
    exit $ret;
fi


# about
echo "install success, " &&
echo " 安装数据库 mysql -uroot -p < docs/db.sql"
echo " you can):     sudo ${INITD} start"
echo "      访问 gin"
echo "      http://127.0.0.1:${PORT}"
echo "bravo Gin root is ${TARGET_INSTALL}"