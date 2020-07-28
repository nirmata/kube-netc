#!/bin/bash

echo "Started clearing old kprobes..."

echo "-:ptcp_get_info" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_v6_connect" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:rtcp_v6_connect" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:rtcp_cleanup_rbuf" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:r__sys_socket" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_retransmit_skb" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_sendmsg" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_cleanup_rbuf" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_close" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:pudp_recvmsg" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:p__sys_bind" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:ptcp_v4_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:pudp_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:pudp_sendmsg" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:rtcp_close" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:rudp_recvmsg" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:r__sys_bind" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:rinet_csk_accept" >> /sys/kernel/debug/tracing/kprobe_events || true
echo "-:p__sys_socket" >> /sys/kernel/debug/tracing/kprobe_events || true

echo "Finished Clearing probes..."


echo "args: $1 $2"
./main $1 $2
