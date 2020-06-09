#!/bin/bash

echo "-:ptcp_get_info" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_v6_connect" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:rtcp_v6_connect" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:rtcp_cleanup_rbuf" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:r__sys_socket" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_retransmit_skb" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_sendmsg" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
ptcp_retransmit_skb
echo "-:ptcp_cleanup_rbuf" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_close" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:pudp_recvmsg" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:p__sys_bind" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:ptcp_v4_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:pudp_destroy_sock" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:pudp_sendmsg" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:rtcp_close" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:rudp_recvmsg" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:r__sys_bind" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:rinet_csk_accept" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
echo "-:p__sys_socket" >> /sys/kernel/debug/tracing/kprobe_events 2>/dev/null || true
__sys_socket
