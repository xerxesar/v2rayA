package iptables

import (
	"fmt"
	"github.com/v2rayA/v2rayA/common/cmds"
	"strings"
)

type tproxy struct {
	watcher *LocalIPWatcher
}

var Tproxy tproxy

func (t *tproxy) AddIPWhitelist(cidr string) {
	// avoid duplication
	t.RemoveIPWhitelist(cidr)
	var commands string
	commands = fmt.Sprintf(`iptables -w 2 -t mangle -I TP_RULE 5 -d %s -j RETURN`, cidr)
	if !strings.Contains(cidr, ".") {
		//ipv6
		commands = strings.Replace(commands, "iptables", "ip6tables", 1)
	}
	cmds.ExecCommands(commands, false)
}

func (t *tproxy) RemoveIPWhitelist(cidr string) {
	var commands string
	commands = fmt.Sprintf(`iptables -w 2 -t mangle -D SETMARK -d %s -j RETURN`, cidr)
	if !strings.Contains(cidr, ".") {
		//ipv6
		commands = strings.Replace(commands, "iptables", "ip6tables", 1)
	}
	cmds.ExecCommands(commands, false)
}

func (t *tproxy) GetSetupCommands() SetupCommands {
	commands := `
ip rule add fwmark 1 table 100
ip route add local 0.0.0.0/0 dev lo table 100

iptables -w 2 -t mangle -N TP_OUT
iptables -w 2 -t mangle -N TP_PRE
iptables -w 2 -t mangle -N TP_RULE

iptables -w 2 -t mangle -I OUTPUT -j TP_OUT
iptables -w 2 -t mangle -I PREROUTING -j TP_PRE

iptables -w 2 -t mangle -A TP_OUT -m mark --mark 0xff -j RETURN
iptables -w 2 -t mangle -A TP_OUT -p tcp -m addrtype --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
iptables -w 2 -t mangle -A TP_OUT -p udp -m addrtype --src-type LOCAL ! --dst-type LOCAL -j TP_RULE

iptables -w 2 -t mangle -A TP_PRE -i lo -m mark ! --mark 1 -j RETURN
iptables -w 2 -t mangle -A TP_PRE -p tcp -m addrtype ! --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
iptables -w 2 -t mangle -A TP_PRE -p udp -m addrtype ! --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
iptables -w 2 -t mangle -A TP_PRE -p tcp -m mark --mark 1 -j TPROXY --on-port 32345 --on-ip 127.0.0.1
iptables -w 2 -t mangle -A TP_PRE -p udp -m mark --mark 1 -j TPROXY --on-port 32345 --on-ip 127.0.0.1

iptables -w 2 -t mangle -A TP_RULE -j CONNMARK --restore-mark
iptables -w 2 -t mangle -A TP_RULE -m mark --mark 1 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -i docker+ -j RETURN
iptables -w 2 -t mangle -A TP_RULE -i veth+ -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 0.0.0.0/32 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 10.0.0.0/8 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 100.64.0.0/10 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 127.0.0.0/8 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 169.254.0.0/16 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 172.16.0.0/12 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 192.0.0.0/24 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 192.0.2.0/24 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 192.88.99.0/24 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 192.168.0.0/16 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 198.18.0.0/15 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 198.51.100.0/24 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 203.0.113.0/24 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 224.0.0.0/4 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -d 240.0.0.0/4 -j RETURN
iptables -w 2 -t mangle -A TP_RULE -p tcp -m tcp --syn -j MARK --set-mark 1
iptables -w 2 -t mangle -A TP_RULE -p udp -m conntrack --ctstate NEW -j MARK --set-mark 1
iptables -w 2 -t mangle -A TP_RULE -j CONNMARK --save-mark
`
	if IsIPv6Supported() {
		commands += `
ip -6 rule add fwmark 1 table 100
ip -6 route add local ::/0 dev lo table 100

ip6tables -w 2 -t mangle -N TP_OUT
ip6tables -w 2 -t mangle -N TP_PRE
ip6tables -w 2 -t mangle -N TP_RULE

ip6tables -w 2 -t mangle -I OUTPUT -j TP_OUT
ip6tables -w 2 -t mangle -I PREROUTING -j TP_PRE

ip6tables -w 2 -t mangle -A TP_OUT -m mark --mark 0xff -j RETURN
ip6tables -w 2 -t mangle -A TP_OUT -p tcp -m addrtype --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
ip6tables -w 2 -t mangle -A TP_OUT -p udp -m addrtype --src-type LOCAL ! --dst-type LOCAL -j TP_RULE

ip6tables -w 2 -t mangle -A TP_PRE -i lo -m mark ! --mark 1 -j RETURN
ip6tables -w 2 -t mangle -A TP_PRE -p tcp -m addrtype ! --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
ip6tables -w 2 -t mangle -A TP_PRE -p udp -m addrtype ! --src-type LOCAL ! --dst-type LOCAL -j TP_RULE
ip6tables -w 2 -t mangle -A TP_PRE -p tcp -m mark --mark 1 -j TPROXY --on-port 32345 --on-ip ::1
ip6tables -w 2 -t mangle -A TP_PRE -p udp -m mark --mark 1 -j TPROXY --on-port 32345 --on-ip ::1

ip6tables -w 2 -t mangle -A TP_RULE -j CONNMARK --restore-mark
ip6tables -w 2 -t mangle -A TP_RULE -m mark --mark 1 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -i docker+ -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -i veth+ -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d ::/128 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d ::1/128 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d 64:ff9b::/96 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d 100::/64 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d 2001::/32 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d 2001:20::/28 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d fe80::/10 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -d ff00::/8 -j RETURN
ip6tables -w 2 -t mangle -A TP_RULE -p tcp -m tcp --syn -j MARK --set-mark 1
ip6tables -w 2 -t mangle -A TP_RULE -p udp -m conntrack --ctstate NEW -j MARK --set-mark 1
ip6tables -w 2 -t mangle -A TP_RULE -j CONNMARK --save-mark
`
	}
	return SetupCommands(commands)
}

func (t *tproxy) GetCleanCommands() CleanCommands {
	commands := `
ip rule del fwmark 1 table 100 
ip route del local 0.0.0.0/0 dev lo table 100

iptables -w 2 -t mangle -F TP_OUT
iptables -w 2 -t mangle -D OUTPUT -j TP_OUT
iptables -w 2 -t mangle -X TP_OUT
iptables -w 2 -t mangle -F TP_PRE
iptables -w 2 -t mangle -D PREROUTING -j TP_PRE
iptables -w 2 -t mangle -X TP_PRE
iptables -w 2 -t mangle -F TP_RULE
iptables -w 2 -t mangle -X TP_RULE
`
	if IsIPv6Supported() {
		commands += `
ip -6 rule del fwmark 1 table 100 
ip -6 route del local ::/0 dev lo table 100

ip6tables -w 2 -t mangle -F TP_OUT
ip6tables -w 2 -t mangle -D OUTPUT -j TP_OUT
ip6tables -w 2 -t mangle -X TP_OUT
ip6tables -w 2 -t mangle -F TP_PRE
ip6tables -w 2 -t mangle -D PREROUTING -j TP_PRE
ip6tables -w 2 -t mangle -X TP_PRE
ip6tables -w 2 -t mangle -F TP_RULE
ip6tables -w 2 -t mangle -X TP_RULE
`
	}
	return CleanCommands(commands)
}
