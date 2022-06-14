#!/bin/bash

##############################################
#user options: ./package.sh -j16
##############################################
jobs=$1; if [[ "" == $jobs ]]; then jobs="-j1";fi


#############################################
# global config
#############################################

ok_msg(){
    echo -e "${1}${POS}${BLACK}[${GREEN} OK ${BLACK}]"
}

failed_msg(){
    echo -e "${1}${POS}${BLACK}[${RED}FAILED${RED}]"
}

#############################################
#discover the current work dir
#############################################

echo "argv[0]=$0"
if [[ ! -f $0 ]]; then
    echo "directly execute the scripts on shell.";
    work_dir=`pwd`
else
    echo "execute scripts in file:$0";
    work_dir=`dirname $0`; work_dir=`(cd ${work_dir} && pwd)`
fi

build_objs="${work_dir}/objs"
package_dir=${build_objs}/package

linux_dist_id=`lsb_release -is 2>/dev/null`
ret=$?; if [[ $ret -ne 0 ]];then
    linux_dist_id=`cat /etc/redhat-release 2>/dev/null|awk '{print $1}'`
fi

linux_dist_version=`lsb_release -rs 2>/dev/null`
ret=$?; if [[ $ret -ne 0 ]];then
    linux_dist_version=`cat /etc/redhat-release 2>/dev/null`
    linux_dist_main_version=`echo $linux_dist_version|awk -F 'release' '{print $2}'|awk '{print $1}'`
else
    linux_dist_main_version=`echo $linux_dist_version`
fi

do_package(){
    branch=`git symbolic-ref --short -q HEAD`
    version=`git rev-parse --short HEAD`
    date_time=`date +%Y_%m_%d_%H_%m`
    package_name="transcoder-${linux_dist_id}${linux_dist_main_version}-${branch}-${version}-${date_time}"
    package_file="${package_name}.zip"
    package_dirname="${package_dir}/${package_name}"
    package_path="${package_dir}/${package_file}"

    echo "start packageing"
    echo "package_file=${package_file}"

    (
        sudo rm -rf ${package_dirname} && mkdir -p ${package_dirname}/objs &&
        cp ${work_dir}/install.sh ${package_dirname}/install.sh &&
        mkdir -p ${package_dirname}/configs &&
        cp -rf ${work_dir}/configs/config.yaml ${package_dirname}/configs/config.yaml &&
        mkdir -p ${package_dirname}/docs
        cp -rf ${work_dir}/docs/db.sql ${package_dirname}/docs/db.sql &&
        cp -rf ${work_dir}/etc ${package_dirname}/ &&
        cp bravo_gin ${package_dirname}/ &&
        mkdir -p ${package_dirname}/logs &&
        cp ${work_dir}/objs/ffmpeg ${package_dirname}/objs/ &&
        cp ${work_dir}/objs/ffprobe ${package_dirname}/objs/ &&
        cp ${work_dir}/objs/generate_hwinfo ${package_dirname}/objs/ &&
        cp ${work_dir}/objs/liblicense-shared.so ${package_dirname}/objs/ &&
        cp ${work_dir}/objs/verify_license ${package_dirname}/objs/ &&
        cd ${package_dir} && rm -f ${package_file} && zip -q -r ${package_file} ${package_name}
    )
    ret=$?; if [[ 0 -ne ${ret} ]]; then failed_msg "package failed"; return $ret; fi

    ok_msg "success: ${package_path}"
}

echo "build core components"
bash build.sh $jobs;
ret=$?; if [[ $ret -ne 0 ]]; then
    failed_msg "build failed."
    exit $ret;
fi

do_package
ret=$?; if [[ $ret -ne 0 ]]; then
    failed_msg "package failed"
    exit $ret;
fi
ok_msg "package success"

echo "install transcoder:"
echo "bash install.sh"