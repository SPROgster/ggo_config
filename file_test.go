package ggo

import "testing"

func (file *Config) checkEntry(t *testing.T, isActive bool, name string, value string, comment string) {
	entry := file.Delete(name)
	if entry == nil {
		t.Errorf("'%s' entry not found", name)
	} else {
		switch ce := entry.(type) {
		case *ConfigEntry:
			if ce.IsActive != isActive {
				t.Errorf("'%s' invalid active state\n", name)
			}
			if ce.Value != value {
				t.Errorf("'%s' invalid value:'%s'\n", name, value)
			}
			if ce.Comment != comment {
				t.Errorf("'%s' invalid active state\n", name)
			}
			// New check place
		default:
			t.Errorf("Invalid type for '%s'\n", name)
		}
	}
}

func (file *Config) checkMultiEntry(t *testing.T, name string, values map[string]bool) {
	entry := file.Delete(name)
	if entry == nil {
		t.Errorf("'%s' entry not found", name)
	} else {
		switch ce := entry.(type) {
		case *ConfigMultiEntry:
			for value, isActive := range values {
				e := ce.Delete(value)
				if e == nil {
					t.Errorf("'%s' multi-entry value '%s' not found", name, value)
				} else {
					if e.Value != value {
						t.Errorf("'%s' multi-entry value '%s' invalid value:'%s'\n", name, value, value)
					}
					if e.IsActive != isActive {
						t.Errorf("'%s' multi-entry value '%s' invalid active state\n", name, value)
					}
				}
			}

		default:
			t.Errorf("Invalid type for '%s'\n", name)
		}
	}
}

func Test_X_GgoConfig_Parse(t *testing.T) {
	testData := []string{
		"sym.prot.ipv4		198.18.1.2/24",
		"sym.prot.vlan		106",
		"sym.raw.ipv4		198.18.0.2/24",
		"sym.raw.vlan		103",
		"",
		"asym.prot.ipv4		198.18.3.2/24",
		"asym.prot.vlan		206",
		"asym.raw.ipv4		198.18.2.2/24",
		"asym.raw.vlan		203",
		"",
		"service.ipv4 198.18.5.2/29",
		"service.vlan 210",
		"",
		"mac		\"ec:93:ed:01:00:00\"",
		"nh.mac		\"78:fe:3d:58:c5:4e\"",
		"",
		"pcap-speed		220",
		"retransmit.skip.net		4.5.6.0/24",
		"",
		"#sflow.drop.pool		0",
		"sflow.drop.rate		0 #1000",
		"sflow.drop.speed		40",
		"#sflow.raw.pool		0",
		"sflow.raw.rate		0 #224",
		"sflow.raw.speed		250",
		"#sflow.accept.pool	0",
		"sflow.accept.rate	0 #1000",
		"sflow.accept.speed	40",
		"",
		"sync	 	  239.0.0.3",
		"sync              239.1.0.3",
		"",
		"#sync-self-session 239.0.0.2",
		"#sync-self-ch-ttl  239.1.0.2",
		"sync-neighbour 198.18.1.3",
		"sync-neighbour 198.18.1.1",
		"sync-neighbour 198.18.1.8",
		"sync-neighbour 198.18.0.8",
		"",
		"pcap-pool 0",
		"#",
		"#sync-neighbour 198.18.1.1",
		"# Bucket configuration",
		"# IPv4",
		"#tb.sym.ipv4_fragmented.32.speed 50",
		"#tb.sym.ipv4_fragmented.bps.24.speed 6250000",
		"#tb.asym.ipv4_fragmented.32.speed 1600",
		"#tb.asym.ipv4_fragmented.bps.24.speed 625000",
		"## TCP",
		"## SYN",
		"#tb.sym.syn.cookie_per_src.32.speed 10",
		"#tb.sym.syn.cookie_per_src.32.setting 1000",
		"#tb.asym.syn.cookie_per_src.32.speed 50",
		"#tb.asym.syn.cookie_per_src.32.setting 0",
		"#tb.sym.syn.ttl.32.speed 1536",
		"#tb.sym.syn.ttl.32.setting 500",
		"#tb.asym.syn.ttl.32.speed 480",
		"#tb.sym.syn.ttl.24.speed 0",
		"#tb.sym.syn.ttl.24.setting 1",
		"#tb.asym.syn.ttl.24.speed 700",
		"#tb.sym.syn.low.32.speed 640",
		"#tb.sym.syn.low.24.speed 20480",
		"#tb.asym.syn.low.32.speed 320",
		"#tb.asym.syn.low.24.speed 4096",
		"#tb.sym.syn.options.32.speed 640",
		"#tb.sym.syn.options.24.speed 20480",
		"#tb.asym.syn.options.32.speed 320",
		"#tb.asym.syn.options.24.speed 5120",
		"#tb.sym.syn.retransmit.32.speed 1600",
		"#tb.sym.syn.retransmit.24.speed 655360",
		"#tb.asym.syn.retransmit.32.speed 1600",
		"#tb.asym.syn.retransmit.24.speed 30720",
		"## SYN+ACK",
		"#tb.asym.syn_ack.low.32.speed 300",
		"#tb.asym.syn_ack.low.24.speed 5120",
		"#tb.asym.syn_ack.final.32.speed 200",
		"#tb.asym.syn_ack.final.24.speed 192000",
		"## SYN_IN_OUT",
		"#tb.sym.tcp.syn.synced.32.setting 1000",
		"#tb.sym.tcp.syn.not_synced.32.speed 1280",
		"#tb.sym.tcp.syn.not_synced.32.setting 500",
		"#tb.sym.tcp.syn.semivalid.24.setting 500",
		"## TCP",
		"#tb.sym.tcp.bad_seq.32.speed 7680",
		"#tb.sym.tcp.bad_seq.32.setting 500",
		"#tb.asym.tcp.bad_seq.32.speed 1280",
		"#tb.sym.tcp.bad_seq.24.speed 0",
		"#tb.sym.tcp.bad_seq.24.setting 500",
		"#tb.sym.tcp.valid.syn.retransmit.32.setting 0",
		"#tb.sym.tcp.semivalid.32.speed 960",
		"#tb.sym.tcp.semivalid.32.setting 500",
		"#tb.sym.tcp.semivalid.24.speed 11200",
		"#tb.sym.tcp.semivalid.24.setting 500",
		"#tb.sym.tcp.semivalid_outgoing.32.speed 0",
		"#tb.sym.tcp.semivalid_outgoing.24.speed 0",
		"#tb.asym.tcp.unknown.32.speed 0",
		"#tb.asym.tcp.unknown.grace.32.speed 0",
		"#tb.asym.tcp.unknown.grace.32.setting 0",
		"#tb.asym.tcp.closing.32.speed 2560",
		"##asymmetric buckets",
		"#tb.sym.memcached_amplifications.32.speed 0",
		"#tb.asym.memcached_amplifications.32.speed 0",
		"#tb.asym.dns.32.speed 7000",
		"#tb.asym.udp.32.speed 20000",
		"#tb.asym.ipv4_others.32.speed 10000",
		"#tb.asym.dns.any.32.speed 2048",
		"#tb.asym.dns.not_any.32.speed 500",
		"#tb.asym.sip.invite.32.speed 1280",
		"#tb.asym.sip.register.32.speed 1280",
		"##symmetric buckets",
		"#tb.sym.udp.32.speed 800",
		"#tb.sym.dns.32.speed 5000",
		"#tb.sym.sip.invite.32.speed 1600",
		"#tb.sym.sip.register.32.speed 1600",
		"#tb.sym.ipv4_others.32.speed 384",
		"## /24",
		"#tb.sym.dns.24.speed 0",
		"#tb.sym.dns.24.setting 1000",
		"#tb.sym.ipv4_others.24.speed 0",
		"#tb.sym.ipv4_others.24.setting 1",
		"#tb.asym.ipv4_others.24.speed 0",
		"#tb.asym.ipv4_others.24.setting 1000",
		"#tb.sym.udp.24.speed 0",
		"#tb.sym.udp.24.setting 1",
		"#tb.asym.syn.low.24.setting 1000",
		"#tb.asym.udp.24.speed 7000",
		"#tb.sym.udp.dst_ports.32.setting 0",
		"#tb.sym.udp.dst_ports.32.speed 800",
		"#tb.sym.udp.32.setting 0",
		"#tb.asym.udp.dst_ports.32.speed 2000",
		"#tb.asym.udp.dst_ports.32.setting 0",
		"#tb.asym.tcp.unknown.24.setting 1000",
		"#tb.asym.tcp.unknown.24.speed 1000",
		"#tb.asym.syn.ttl.24.setting 1000",
		"#tb.asym.ntp_amplifications.32.speed 640",
		"#tb.asym.ntp_amplifications.32.setting 1000",
		"#tb.sym.tcp.semivalid_outgoing.32.setting 1",
		"#tb.sym.tcp.semivalid_outgoing.24.setting 1",
		"#tb.asym.memcached_amplifications.32.setting 0",
		"#tb.sym.memcached_amplifications.32.setting 0",
		"#tb.asym.syn.ttl.32.setting 500",
		"#tb.asym.upnp_amplifications.24.speed 10000",
		"#tb.asym.syn.retransmit.24.setting 1",
		"#tb.asym.syn.retransmit.32.setting 1000",
		"#tb.sym.ipv4_fragmented.32.setting 0",
		"#tb.asym.teamspeak.32.setting 250",
		"#tb.asym.teamspeak.32.speed 500",
		"#tb.asym.teamspeak.24.speed 10240",
		"#tb.asym.teamspeak.24.setting 1000",
		"#tb.sym.syn.options.32.setting 1000",
		"#tb.sym.syn.retransmit.32.setting 1000",
		"#tb.sym.syn.retransmit.24.setting 1000",
		"#tb.sym.syn.low.32.setting 1000",
		"#tb.asym.unreal_tournament.32.speed 0",
		"#tb.asym.udp.32.setting 1000",
		"#tb.asym.udp.session_add.32.speed 4096",
		"#tb.asym.udp.session_add.32.setting 0",
		"#tb.asym.udp.24.setting 0",
		"#tb.asym.syn_ack.final.32.setting 0",
		"#tb.asym.quic_like.32.setting 100",
		"#tb.asym.quic_like.32.speed 2000",
		"#tb.asym.quic_like.24.speed 3000",
		"#tb.asym.quic_like.24.setting 100",
		"#tb.sym.teamspeak.32.setting 0",
		"#tb.sym.teamspeak.32.speed 320",
		"#tb.sym.teamspeak.24.speed 10000",
		"#tb.sym.teamspeak.24.setting 0",
		"eth-0_1",
		"eth-0_0",
		"eth-1_3",
		"eth-1_0",
		"eth-2_3",
		"eth-2_2",
		"eth-3_0",
		"#cores-per-port 8",
		"#switch off cookie filter",
		"",
		"",
		"#tb.sym.syn.cookie_per_src.32.speed           0",
		"#tb.sym.syn.cookie_per_src.32.setting         1",
		"",
		"#tb.sym.syn.ttl2sync.32.setting                   1",
		"#tb.asym.syn.ttl2sync.32.setting                  1",
		"",
		"#tb.sym.syn.ttl2ch.32.setting                        1",
		"#tb.asym.syn.ttl2ch.32.setting                       1",
		"",
		"#tb.sym.syn.ch2sync.32.setting                       1",
		"#tb.asym.syn.ch2sync.32.setting                      1",
		"",
		"#tb.sym.syn.ch2ttl.32.setting                        1",
		"#tb.asym.syn.ch2ttl.32.setting                       1",
		"",
		"#tb.sym.syn.sync2ttl.32.setting                      1",
		"#tb.asym.syn.sync2ttl.32.setting                     1",
		"",
		"#tb.sym.syn.sync2ch.32.setting                       1",
		"#tb.asym.syn.sync2ch.32.setting                      1",
		"",
		"#tb.sym.syn.low.32.setting 0",
		"#tb.sym.syn.low.32.speed 0",
		"#tb.sym.syn.low.24.setting 0",
		"#tb.sym.syn.low.24.speed 0",
		"tb.sym.syn.options.24.setting           0",
		"tb.sym.syn.options.24.speed             0",
		"tb.sym.syn.options.32.setting           0",
		"tb.sym.syn.options.32.speed             0",
		"#tb.sym.syn.retransmit.24.setting                0",
		"#tb.sym.syn.retransmit.24.speed          0",
		"#tb.sym.syn.retransmit.32.setting                1",
		"#tb.sym.syn.retransmit.32.speed          0",
		"",
	}

	file := NewConfig()
	file.SetKeyMultiple("sync", true)
	file.SetKeyMultiple("sync-neighbour", true)
	file.FromStrings(testData)

	file.checkEntry(t, true, "sym.prot.ipv4", "198.18.1.2/24", "")
	file.checkEntry(t, true, "sym.prot.vlan", "106", "")
	file.checkEntry(t, true, "sym.raw.ipv4", "198.18.0.2/24", "")
	file.checkEntry(t, true, "sym.raw.vlan", "103", "")
	file.checkEntry(t, true, "asym.prot.ipv4", "198.18.3.2/24", "")
	file.checkEntry(t, true, "asym.prot.vlan", "206", "")
	file.checkEntry(t, true, "asym.raw.ipv4", "198.18.2.2/24", "")
	file.checkEntry(t, true, "asym.raw.vlan", "203", "")
	file.checkEntry(t, true, "service.ipv4", "198.18.5.2/29", "")
	file.checkEntry(t, true, "service.vlan", "210", "")
	file.checkEntry(t, true, "mac", "\"ec:93:ed:01:00:00\"", "")
	file.checkEntry(t, true, "nh.mac", "\"78:fe:3d:58:c5:4e\"", "")
	file.checkEntry(t, true, "pcap-speed", "220", "")
	file.checkEntry(t, true, "retransmit.skip.net", "4.5.6.0/24", "")
	file.checkEntry(t, false, "sflow.drop.pool", "0", "")
	file.checkEntry(t, true, "sflow.drop.rate", "0", "1000")
	file.checkEntry(t, true, "sflow.drop.speed", "40", "")
	file.checkEntry(t, false, "sflow.raw.pool", "0", "")
	file.checkEntry(t, true, "sflow.raw.rate", "0", "224")
	file.checkEntry(t, true, "sflow.raw.speed", "250", "")
	file.checkEntry(t, false, "sflow.accept.pool", "0", "")
	file.checkEntry(t, true, "sflow.accept.rate", "0", "1000")
	file.checkEntry(t, true, "sflow.accept.speed", "40", "")
	file.checkEntry(t, false, "sync-self-session", "239.0.0.2", "")
	file.checkEntry(t, false, "sync-self-ch-ttl", "239.1.0.2", "")
	file.checkEntry(t, true, "pcap-pool", "0", "")
	file.checkEntry(t, false, "Bucket", "configuration", "")
	file.checkEntry(t, false, "IPv4", "", "")
	file.checkEntry(t, false, "tb.sym.ipv4_fragmented.32.speed", "50", "")
	file.checkEntry(t, false, "tb.sym.ipv4_fragmented.bps.24.speed", "6250000", "")
	file.checkEntry(t, false, "tb.asym.ipv4_fragmented.32.speed", "1600", "")
	file.checkEntry(t, false, "tb.asym.ipv4_fragmented.bps.24.speed", "625000", "")
	file.checkEntry(t, false, "TCP", "", "")
	file.checkEntry(t, false, "SYN", "", "")
	file.checkEntry(t, false, "tb.asym.syn.cookie_per_src.32.speed", "50", "")
	file.checkEntry(t, false, "tb.asym.syn.cookie_per_src.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl.32.speed", "1536", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl.32.setting", "500", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl.32.speed", "480", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl.24.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl.24.speed", "700", "")
	file.checkEntry(t, false, "tb.asym.syn.low.32.speed", "320", "")
	file.checkEntry(t, false, "tb.asym.syn.low.24.speed", "4096", "")
	file.checkEntry(t, false, "tb.asym.syn.options.32.speed", "320", "")
	file.checkEntry(t, false, "tb.asym.syn.options.24.speed", "5120", "")
	file.checkEntry(t, false, "tb.asym.syn.retransmit.32.speed", "1600", "")
	file.checkEntry(t, false, "tb.asym.syn.retransmit.24.speed", "30720", "")
	file.checkEntry(t, false, "SYN+ACK", "", "")
	file.checkEntry(t, false, "tb.asym.syn_ack.low.32.speed", "300", "")
	file.checkEntry(t, false, "tb.asym.syn_ack.low.24.speed", "5120", "")
	file.checkEntry(t, false, "tb.asym.syn_ack.final.32.speed", "200", "")
	file.checkEntry(t, false, "tb.asym.syn_ack.final.24.speed", "192000", "")
	file.checkEntry(t, false, "SYN_IN_OUT", "", "")
	file.checkEntry(t, false, "tb.sym.tcp.syn.synced.32.setting", "1000", "")
	file.checkEntry(t, false, "tb.sym.tcp.syn.not_synced.32.speed", "1280", "")
	file.checkEntry(t, false, "tb.sym.tcp.syn.not_synced.32.setting", "500", "")
	file.checkEntry(t, false, "tb.sym.tcp.syn.semivalid.24.setting", "500", "")
	file.checkEntry(t, false, "tb.sym.tcp.bad_seq.32.speed", "7680", "")
	file.checkEntry(t, false, "tb.sym.tcp.bad_seq.32.setting", "500", "")
	file.checkEntry(t, false, "tb.asym.tcp.bad_seq.32.speed", "1280", "")
	file.checkEntry(t, false, "tb.sym.tcp.bad_seq.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.tcp.bad_seq.24.setting", "500", "")
	file.checkEntry(t, false, "tb.sym.tcp.valid.syn.retransmit.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid.32.speed", "960", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid.32.setting", "500", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid.24.speed", "11200", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid.24.setting", "500", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid_outgoing.32.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid_outgoing.24.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.tcp.unknown.32.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.tcp.unknown.grace.32.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.tcp.unknown.grace.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.tcp.closing.32.speed", "2560", "")
	file.checkEntry(t, false, "asymmetric", "buckets", "")
	file.checkEntry(t, false, "tb.sym.memcached_amplifications.32.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.memcached_amplifications.32.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.dns.32.speed", "7000", "")
	file.checkEntry(t, false, "tb.asym.udp.32.speed", "20000", "")
	file.checkEntry(t, false, "tb.asym.ipv4_others.32.speed", "10000", "")
	file.checkEntry(t, false, "tb.asym.dns.any.32.speed", "2048", "")
	file.checkEntry(t, false, "tb.asym.dns.not_any.32.speed", "500", "")
	file.checkEntry(t, false, "tb.asym.sip.invite.32.speed", "1280", "")
	file.checkEntry(t, false, "tb.asym.sip.register.32.speed", "1280", "")
	file.checkEntry(t, false, "symmetric", "buckets", "")
	file.checkEntry(t, false, "tb.sym.udp.32.speed", "800", "")
	file.checkEntry(t, false, "tb.sym.dns.32.speed", "5000", "")
	file.checkEntry(t, false, "tb.sym.sip.invite.32.speed", "1600", "")
	file.checkEntry(t, false, "tb.sym.sip.register.32.speed", "1600", "")
	file.checkEntry(t, false, "tb.sym.ipv4_others.32.speed", "384", "")
	file.checkEntry(t, false, "/24", "", "")
	file.checkEntry(t, false, "tb.sym.dns.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.dns.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.sym.ipv4_others.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.ipv4_others.24.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.ipv4_others.24.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.ipv4_others.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.sym.udp.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.udp.24.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.low.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.asym.udp.24.speed", "7000", "")
	file.checkEntry(t, false, "tb.sym.udp.dst_ports.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.udp.dst_ports.32.speed", "800", "")
	file.checkEntry(t, false, "tb.sym.udp.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.udp.dst_ports.32.speed", "2000", "")
	file.checkEntry(t, false, "tb.asym.udp.dst_ports.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.tcp.unknown.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.asym.tcp.unknown.24.speed", "1000", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.asym.ntp_amplifications.32.speed", "640", "")
	file.checkEntry(t, false, "tb.asym.ntp_amplifications.32.setting", "1000", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid_outgoing.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.tcp.semivalid_outgoing.24.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.memcached_amplifications.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.memcached_amplifications.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl.32.setting", "500", "")
	file.checkEntry(t, false, "tb.asym.upnp_amplifications.24.speed", "10000", "")
	file.checkEntry(t, false, "tb.asym.syn.retransmit.24.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.retransmit.32.setting", "1000", "")
	file.checkEntry(t, false, "tb.sym.ipv4_fragmented.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.teamspeak.32.setting", "250", "")
	file.checkEntry(t, false, "tb.asym.teamspeak.32.speed", "500", "")
	file.checkEntry(t, false, "tb.asym.teamspeak.24.speed", "10240", "")
	file.checkEntry(t, false, "tb.asym.teamspeak.24.setting", "1000", "")
	file.checkEntry(t, false, "tb.asym.unreal_tournament.32.speed", "0", "")
	file.checkEntry(t, false, "tb.asym.udp.32.setting", "1000", "")
	file.checkEntry(t, false, "tb.asym.udp.session_add.32.speed", "4096", "")
	file.checkEntry(t, false, "tb.asym.udp.session_add.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.udp.24.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.syn_ack.final.32.setting", "0", "")
	file.checkEntry(t, false, "tb.asym.quic_like.32.setting", "100", "")
	file.checkEntry(t, false, "tb.asym.quic_like.32.speed", "2000", "")
	file.checkEntry(t, false, "tb.asym.quic_like.24.speed", "3000", "")
	file.checkEntry(t, false, "tb.asym.quic_like.24.setting", "100", "")
	file.checkEntry(t, false, "tb.sym.teamspeak.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.teamspeak.32.speed", "320", "")
	file.checkEntry(t, false, "tb.sym.teamspeak.24.speed", "10000", "")
	file.checkEntry(t, false, "tb.sym.teamspeak.24.setting", "0", "")
	file.checkEntry(t, true, "eth-0_1", "", "")
	file.checkEntry(t, true, "eth-0_0", "", "")
	file.checkEntry(t, true, "eth-1_3", "", "")
	file.checkEntry(t, true, "eth-1_0", "", "")
	file.checkEntry(t, true, "eth-2_3", "", "")
	file.checkEntry(t, true, "eth-2_2", "", "")
	file.checkEntry(t, true, "eth-3_0", "", "")
	file.checkEntry(t, false, "cores-per-port", "8", "")
	file.checkEntry(t, false, "tb.sym.syn.cookie_per_src.32.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.cookie_per_src.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl2sync.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl2sync.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.ttl2ch.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.ttl2ch.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.ch2sync.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.ch2sync.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.ch2ttl.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.ch2ttl.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.sync2ttl.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.sync2ttl.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.sync2ch.32.setting", "1", "")
	file.checkEntry(t, false, "tb.asym.syn.sync2ch.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.low.32.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.low.32.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.low.24.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.low.24.speed", "0", "")
	file.checkEntry(t, true, "tb.sym.syn.options.24.setting", "0", "")
	file.checkEntry(t, true, "tb.sym.syn.options.24.speed", "0", "")
	file.checkEntry(t, true, "tb.sym.syn.options.32.setting", "0", "")
	file.checkEntry(t, true, "tb.sym.syn.options.32.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.retransmit.24.setting", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.retransmit.24.speed", "0", "")
	file.checkEntry(t, false, "tb.sym.syn.retransmit.32.setting", "1", "")
	file.checkEntry(t, false, "tb.sym.syn.retransmit.32.speed", "0", "")

	values := map[string]bool{
		"239.0.0.3": true,
		"239.1.0.3": true,
	}
	file.checkMultiEntry(t, "sync", values)

	values = map[string]bool{
		"198.18.1.3": true,
		"198.18.1.1": true,
		"198.18.1.8": true,
		"198.18.0.8": true,
	}
	file.checkMultiEntry(t, "sync-neighbour", values)

	if len(file.fields) != 0 {
		t.Errorf("Some fields (%d) left unprocessed %v\n", len(file.fields), file.fields)
	}
}

func Test_GgoConfig_Merge(t *testing.T) {
	testData := [][]string{
		{
			"sync	 	  239.0.0.3",
			"sync              239.1.0.3",
			"",
			"#tb.sym.ipv4_fragmented.32.speed 50",
			"#tb.sym.ipv4_fragmented.bps.24.speed 6250000",
			"#tb.asym.ipv4_fragmented.32.speed 1600",
			"#tb.asym.ipv4_fragmented.bps.24.speed 625000",
		},
		{
			"#sync	 	  239.0.0.3",	// Тут деактивировали
			"sync              239.1.0.3",
			"sync              239.2.0.3", // Это добавили
			"",
			"#tb.sym.ipv4_fragmented.32.speed 50",
			"tb.sym.ipv4_fragmented.bps.24.speed 1234", // Тут изменили и активировали
			"#tb.asym.ipv4_fragmented.32.speed 1600",
			"#tb.asym.ipv4_fragmented.bps.24.speed 625000 # some comment", // Сюда добавили комментарий
		},
	}

	file1 := NewConfig()
	file1.SetKeyMultiple("sync", true)
	file1.SetKeyMultiple("sync-neighbour", true)
	file1.FromStrings(testData[0])

	file2 := file1.CopyScheme()
	file2.FromStrings(testData[1])

	result := Merge(file1, file2)
	values := map[string]bool{
		"239.0.0.3": false,
		"239.1.0.3": true,
		"239.2.0.3": true,
	}

	result.checkMultiEntry(t, "sync", values)

	result.checkEntry(t, false, "tb.sym.ipv4_fragmented.32.speed", "50", "")
	result.checkEntry(t, true, "tb.sym.ipv4_fragmented.bps.24.speed", "1234", "")
	result.checkEntry(t, false, "tb.asym.ipv4_fragmented.32.speed", "1600", "")
	result.checkEntry(t, false, "tb.asym.ipv4_fragmented.bps.24.speed", "625000", "some comment")


	if len(result.fields) != 0 {
		t.Errorf("Some fields (%d) left unprocessed %v\n", len(result.fields), result.fields)
	}
}

func Test_GgoConfig_MergeWithEmpties(t *testing.T) {
	testData := [][]string{
		{
			"sync	 	  239.0.0.3",
			"sync              239.1.0.3",
			"",
			"#tb.sym.ipv4_fragmented.32.speed 50",
			"#tb.sym.ipv4_fragmented.bps.24.speed 6250000",
			"#tb.asym.ipv4_fragmented.32.speed 1600",
			"#tb.asym.ipv4_fragmented.bps.24.speed 625000",
		},
		{
			"#sync	 	  239.0.0.3",	// Тут деактивировали
			"sync              239.1.0.3",
			"sync              239.2.0.3", // Это добавили
			"",
			"#tb.sym.ipv4_fragmented.32.speed 50",
			"tb.sym.ipv4_fragmented.bps.24.speed 1234", // Тут изменили и активировали
			"#tb.asym.ipv4_fragmented.32.speed 1600",
			"#tb.asym.ipv4_fragmented.bps.24.speed 625000 # some comment", // Сюда добавили комментарий
		},
	}

	file1 := NewConfig()
	file1.SetKeyMultiple("sync", true)
	file1.SetKeyMultiple("sync-neighbour", true)
	file1.FromStrings(testData[0])

	file2 := NewConfig()

	res := Merge(file1, file2, nil)
	if res == nil {
		t.Error("Tail empty failed")
	}

	res = Merge(file1, nil, file2)
	if res == nil {
		t.Error("Middle empty failed")
	}

	res = Merge(nil, file1, file2)
	if res == nil {
		t.Error("Head empty failed")
	}

	res = Merge(nil, nil, nil)
	if res != nil {
		t.Error("Full empty failed")
	}
}