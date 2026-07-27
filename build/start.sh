#!/bin/bash
BASE_DIR=$(cd `dirname $0`; pwd)
ROOT_DIR=${BASE_DIR}/..
DST_PATH=/opt/csp
#base log define
log_cmd_info="eval echo \"[$0]\" @\`date +\"%Y%m%d %T\"\` [info]:"
log_cmd_error="eval echo \"[$0]\" @\`date +\"%Y%m%d %T\"\` [error]:"
optemg="/optemg"

# ban core dump
ulimit -c  0
#service name define
export SERVICE_NAME=gids
export MALLOC_ARENA_MAX=1
#check process
export CIPHER_PATH=/opt/csp/gids/cipher
pid_service=$(ps -ef|grep gids|grep /opt/csp/gids/module|grep -v grep|awk '{print $2}')

innerIP=""
if [ -n "$FABRIC_ETH" ]; then
    innerIP=$(ip -4 addr show ${FABRIC_ETH} | grep inet | awk '{print $2}' | cut -d'/' -f1)
    echo "FABRIC_ETH IP: $innerIP"
else
    innerIP=$(ip -4 addr show bond-base | grep inet | awk '{print $2}' | cut -d'/' -f1)
    echo "bond-base IP: $innerIP"
fi

if [ -z "${innerIP}" ]; then
    echo "failed to get ip neither ${FABRIC_ETH} nor bond-base"
    exit 1
fi

# csp framework will use FABRIC_IP as advertise ip(use MBASE_ETH if no set FABRIC_IP).
# NEO have used bond-external as MBASE_ETH,so we set bond-base as FABRIC_IP for internal communicating
export FABRIC_IP=${innerIP}

#modify chassis.yaml
id=${APPID}
if [  ! -z "$id" ];then
   ${log_cmd_info} "APPID is: ${APPID}."
    sed -i -e 's/APPLICATION_ID: CSP/APPLICATION_ID: '${APPID}'/' ${DST_PATH}/${SERVICE_NAME}/module/conf/chassis.yaml
    sed -i -e 's/logger_file: \/opt\/csplog\/gids\/gids.log/logger_file: \/opt\/csplog\/'${APPID}'\/'${SERVICE_NAME}'\/gids.log/' ${DST_PATH}/${SERVICE_NAME}/module/conf/lager.yaml

    sed  -i -e 's/^appid.*$/appid = '${APPID}'/g' ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
fi

#modify app.conf
sed -i "s@{config_url}@${MUEN_CONFIG_URL}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
sed -i "s@{tiktok_endpoint}@${TIKTOK_ENDPOINT}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
sed -i "s@{port}@${PORT}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
sed -i "s@{tls_port}@${TLS_PORT}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
sed -i "s@{httpsconfig_url}@${HTTPS_MUEN_CONFIG_URL}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf
sed -i "s@{httpstiktok_endpoint}@${HTTPS_TIKTOK_ENDPOINT}@g" ${DST_PATH}/${SERVICE_NAME}/module/conf/app.conf


if [ -z "$APPID" ];then
        export APPID=0
fi

#log msg
export LOG_PATH="/opt/csplog/${APPID}/${SERVICE_NAME}"
if [ ! -d "$LOG_PATH" ]; then
    mkdir -p $LOG_PATH
fi

#valid ip
function valid_ip()
{
    local ip=${1}
    VALID_CHECK=$(echo ${ip} | awk -F . '$1<=255&&$2<=255&&$3<=255&&$4<=255{print "yes"}')
    varIp=$(echo ${ip}|grep -E "^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$" | xargs echo)
    if [ -n "${varIp}" ]; then
        if [ "${VALID_CHECK}" == "yes" ]; then
            echo "ip $ip available."
            return 0
        else
            echo "ip $ip not available!"
            return 1
        fi
    else
        echo "ip format error!"
        return 1
    fi
}

#get ip
function getLocalIP(){
    filename=/opt/csp_service_ip.ini
    if [ -f "$filename" ];then
        LOCAL_IP=$(grep $SERVICE_NAME /opt/csp_service_ip.ini| awk -F':' '{print $2}' | sed  's/ //g')
        if valid_ip $LOCAL_IP;then
            export SERVICE_IP="${LOCAL_IP}"
            return
        fi
    fi
    if [ x"" != x${MBASE_ETH} ];then
        LOCAL_IP=`ifconfig ${MBASE_ETH} | grep 'inet' | grep -v 'inet6' | awk '{print $2}' | sed 's/addr://g'`
        if valid_ip $LOCAL_IP;then
            return
        fi
    fi
    LOCAL_IP=`ps -ef | grep '/usr/local/bin/kubelet' | grep -v grep|awk -F '--node-ip=' '{print $2}'|awk -F ' ' '{print $1}'`
    if valid_ip $LOCAL_IP;then
        return
    fi
    LOCAL_IP=`ifconfig eth0 | grep 'inet'|grep -v 'inet6'| awk '{print $2}' |sed 's/addr://g'`
    if valid_ip $LOCAL_IP;then
        return
    fi
}
getLocalIP

sudo bash ${DST_PATH}/${SERVICE_NAME}/module/init.sh ${APPID}
#export env
if [ -d "/opt/cspcertssl" ] || [ "${INNER_TLS_MODE}" = "k8s" ]; then
    export CSP_NEW_INNER_CERT_PATH=/opt/cspcertssl
elif [ -d "/optemg/etc/common" ]; then
    export CSP_NEW_INNER_CERT_PATH=/optemg/etc/common/cert
else
    export CSP_NEW_INNER_CERT_PATH=/opt/container/envinfo/cert
fi
PAAS_CRYPTO_PATH=${CIPHER_PATH}/
AGENT_CERT_PATH=/opt/csp/gids/cert/

export SSLPATH="${AGENT_CERT_PATH}"

export CIPHER_ROOT=${AGENT_CERT_PATH}
if [ "X${INSTALL_MODE}" != "X" ] && ([ ${INSTALL_MODE} = "build" ] || [ ${INSTALL_MODE} = "level_nfc" ] || [ ${INSTALL_MODE} = "level" ]) ;then
   $log_cmd_info "Set cipher root to paas."
    export CIPHER_ROOT=${PAAS_CRYPTO_PATH}
fi

export GENERAL_SCRIPT_PATH=/opt/csp/scripts
export CHASSIS_HOME=${DST_PATH}/${SERVICE_NAME}/module
if [ ${EDGE_MODE} != "edge" ];then
    export SSL_CERT_PATH=${AGENT_CERT_PATH}
fi
export AGENT_MANAGER_DIR=/opt/csp/gids

certDir=(
/opt/csp/cert/gids
/opt/paas/srv
)

excludeArray=(
/opt/csplog
/opt/csp/gids/module/gids
/opt/csp/gids/scripts
)

chmod 750 ${LOG_PATH}

chmod 600 -R /opt/csplog/${APPID}/gids/*.log
chmod 600 -R /opt/csplog/${APPID}/gids/*.gz
chmod 640 -Rh /opt/csplog/${APPID}/gids/RScheck/*.log
chmod 440 -Rh /opt/csplog/${APPID}/gids/RScheck/*.zip
chmod 700 -R /opt/csp/default/*
chmod 600 -R /opt/csp/default/cert/*
chmod 700 -R /opt/csp/gids/*
chmod 600 -R /opt/csp/gids/cert/*
chmod 600 -R ${CIPHER_PATH}/*
chmod 600 -R /opt/csp/default/cipher/*
chmod 500 /opt/csp/default/cipher/updateKey.sh
chmod 500 ${CIPHER_PATH}/updateKey.sh

export GRPC_GO_RETRY="on"
export SERVICENAME="gids"
RScli set process /opt/csp/gids/module/gids fd -status true -threshold 3000 -killOnError false thread -status true -threshold 2000 -killOnError false
RScli set process /opt/csp/gids/module/gids fd -alarmStatus on -alarm 95 -clear 90  -alarmCheckTimes 10 -alarmClearTimes 10 -period 30s thread -alarmStatus on -alarm 95 -clear 90 -alarmCheckTimes 10 -alarmClearTimes 10 -period 30s
RScli set container cpu -status true -period 60s -initialDelay 30m -alarmStatus on -alarm 95 -clear 90 -alarmCheckTimes 10 -alarmClearTimes 10 mem -status true -period 60s -initialDelay 3m -alarmStatus on -alarm 95 -clear 90 -alarmCheckTimes 10 -alarmClearTimes 10
RScli set startShell /opt/csp/gids/scripts/start.sh
RScli reload


# update cert
rm -rf /opt/csp/gids/cert/revoke.crl
rm -rf /opt/csp/gids/cert/p12.cert
rm -rf /opt/csp/gids/cert/server.p12
/opt/csp/cspcerttools/cspcerttool -sslPath /opt/csp/gids/cert -cipherPath ${CIPHER_PATH}  -cspCertSslPath ${CSP_NEW_INNER_CERT_PATH}

#export env 边缘场景中心场景默认设置为8
# 防止边缘场景下根据cpu核心数设置GOMAXPROCS导致拉起过多线程, 引起无用的rss占用和性能损耗
export  GOMAXPROCS=8
export  GOGC=50

# Compatible Scenarios
if [[ "${ENABLE_CHAIN}X" == "X" ]]; then
  export ENABLE_HTTP="true"
fi

#execute process
cd ${DST_PATH}/${SERVICE_NAME}/module
command="${DST_PATH}/${SERVICE_NAME}/module/gids -node-endpoint ${GIDSExtendEndpoint} -node-httpsendpoint ${GIDSHttpsExtendEndpoint} >> /opt/csplog/0/gids/gids.log 2>&1 &"

USER="$(whoami)"
if [ "root" = "${USER}" ]; then
    su paas -c "${command}"
else
   ${command}
fi

certDir=(
)

excludeArray=(
)


#check process
pid_service=$(ps -ef|grep gids|grep /opt/csp/gids/module|grep -v grep|awk '{print $2}')

if [ ! -z "${pid_service}" ]; then
    ${log_cmd_info} "gids is alreadly started!"
    exit 0
fi
${log_cmd_info} "end to exec $0."