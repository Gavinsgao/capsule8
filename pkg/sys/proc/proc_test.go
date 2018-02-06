// Copyright 2017 Capsule8, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proc

import (
	"bytes"
	"fmt"
	"testing"
)

// This data was taken from actual /proc/PID/stat files.
var statTests = []struct {
	statFile   string
	pid        int
	comm       string
	ppid       int
	startTime  uint64
	startStack uint64
}{
	{
		statFile:   "4018 (bash) S 4011 4018 4018 34834 7516 4194304 8082 41779 1 85 33 7 115 329 20 0 1 0 8810 24444928 1667 18446744073709551615 4194304 5192876 140734725904528 140734725903192 140515966087290 0 65536 3670020 1266777851 1 0 0 17 3 0 0 1 0 0 7290352 7326856 31535104 140734725912793 140734725912798 140734725912798 140734725914606 0\n",
		pid:        4018,
		comm:       "bash",
		ppid:       4011,
		startTime:  8810,
		startStack: 140734725904528,
	},
	{
		statFile:   "899 (rs:main Q:Reg) S 1 828 828 0 -1 1077936192 720 0 5 0 374 450 0 0 20 0 4 0 512 262553600 2530 18446744073709551615 1 1 0 0 0 0 2146172671 16781830 1132545 0 0 0 -1 1 0 0 9 0 0 0 0 0 0 0 0 0 0\n",
		pid:        899,
		comm:       "rs:main Q:Reg",
		ppid:       1,
		startTime:  512,
		startStack: 0,
	},
	{
		statFile:   "25663 (a b) S 4090 25663 4090 34833 25831 4194304 112 0 0 0 0 0 0 0 20 0 1 0 2591294 4616192 191 18446744073709551615 93931362930688 93931363074588 140721799437360 140721799436056 139765395259690 0 0 0 65538 1 0 0 17 0 0 0 0 0 0 93931365175144 93931365179936 93931378774016 140721799438724 140721799438752 140721799438752 140721799442404 0\n",
		pid:        25663,
		comm:       "a b",
		ppid:       4090,
		startTime:  2591294,
		startStack: 140721799437360,
	},
	{
		statFile:   "25666 ((c) S 4090 25666 4090 34833 25831 4194304 111 0 0 0 0 0 0 0 20 0 1 0 2591294 4616192 197 18446744073709551615 94586441084928 94586441228828 140737160769408 140737160768104 140708343980330 0 0 0 65538 1 0 0 17 3 0 0 0 0 0 94586443329384 94586443334176 94586462375936 140737160774023 140737160774050 140737160774050 140737160777701 0\n",
		pid:        25666,
		comm:       "(c",
		ppid:       4090,
		startTime:  2591294,
		startStack: 140737160769408,
	},
	{
		statFile:   "25669 (d)) S 4090 25669 4090 34833 25831 4194304 114 0 0 0 0 0 0 0 20 0 1 0 2591295 4616192 201 18446744073709551615 93918460887040 93918461030940 140727364187808 140727364186504 140658074984746 0 0 0 65538 1 0 0 17 3 0 0 0 0 0 93918463131496 93918463136288 93918473555968 140727364190599 140727364190626 140727364190626 140727364194277 0\n",
		pid:        25669,
		comm:       "d)",
		ppid:       4090,
		startTime:  2591295,
		startStack: 140727364187808,
	},
	{
		statFile:   "25672 (((e))) S 4090 25672 4090 34833 25831 4194304 114 0 0 0 0 0 0 0 20 0 1 0 2591295 4616192 178 18446744073709551615 94113212719104 94113212863004 140724070346384 140724070345080 140031172235562 0 0 0 65538 1 0 0 17 0 0 0 0 0 0 94113214963560 94113214968352 94113226104832 140724070355326 140724070355356 140724070355356 140724070359010 0\n",
		pid:        25672,
		comm:       "((e))",
		ppid:       4090,
		startTime:  2591295,
		startStack: 140724070346384,
	},
	{
		statFile:   "25675 ( f  ) S 4090 25675 4090 34833 25831 4194304 111 0 0 0 0 0 0 0 20 0 1 0 2591295 4616192 191 18446744073709551615 94829034725376 94829034869276 140737237421792 140737237420488 139937926709546 0 0 0 65538 1 0 0 17 2 0 0 0 0 0 94829036969832 94829036974624 94829068091392 140737237426561 140737237426590 140737237426590 140737237430243 0\n",
		pid:        25675,
		comm:       " f  ",
		ppid:       4090,
		startTime:  2591295,
		startStack: 140737237421792,
	},
}

// TestStatParse tests the parsing of /proc/PID/stat files.
func TestStatParse(t *testing.T) {
	for _, td := range statTests {
		ps := &ProcessStatus{statFields: statFields(td.statFile)}

		if ps.PID() != td.pid {
			t.Errorf("For proc.(*ProcessStatus)PID(), want %d, got %d\n", td.pid, ps.PID())
		}
		if ps.Command() != td.comm {
			t.Errorf("For proc.(*ProcessStatus)Command(), want \"%s\", got \"%s\"\n", td.comm, ps.Command())
		}
		if ps.ParentPID() != td.ppid {
			t.Errorf("For proc.(*ProcessStatus)ParentPID(), want %d, got %d\n", td.ppid, ps.ParentPID())
		}
		if ps.StartTime() != td.startTime {
			t.Errorf("For proc.(*ProcessStatus)startTime(), want %d, got %d\n", td.startTime, ps.StartTime())
		}
		if ps.StartStack() != td.startStack {
			t.Errorf("For proc.(*ProcessStatus)StartStack(), want %d, got %d\n", td.startStack, ps.StartStack())
		}
	}
}

var statusTests = []struct {
	statusFile string
	TGID, PID  int
	UID        []int
	GID        []int
	Name       string
}{
	{
		statusFile: `Name:	systemd
State:	S (sleeping)
Tgid:	1
Ngid:	0
Pid:	1
PPid:	0
TracerPid:	0
Uid:	0	0	0	0
Gid:	0	0	0	0
FDSize:	64
Groups:
NStgid:	1
NSpid:	1
NSpgid:	1
NSsid:	1
VmPeak:	   38928 kB
VmSize:	   37844 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	    5972 kB
VmRSS:	    5832 kB
VmData:	    1596 kB
VmStk:	     132 kB
VmExe:	    1392 kB
VmLib:	    3664 kB
VmPTE:	      92 kB
VmPMD:	      12 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
Threads:	1
SigQ:	0/15593
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	7be3c0fe28014a03
SigIgn:	0000000000001000
SigCgt:	00000001800004ec
CapInh:	0000000000000000
CapPrm:	0000003fffffffff
CapEff:	0000003fffffffff
CapBnd:	0000003fffffffff
CapAmb:	0000000000000000
Seccomp:	0
Cpus_allowed:	ffffffff,ffffffff,ffffffff,ffffffff
Cpus_allowed_list:	0-127
Mems_allowed:	00000000,00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	10752
nonvoluntary_ctxt_switches:	657`,
		TGID: 1,
		PID:  1,
		UID:  []int{0, 0, 0, 0},
		GID:  []int{0, 0, 0, 0},
		Name: "systemd",
	},
	{
		statusFile: `Name:	vmhgfs-fuse
State:	S (sleeping)
Tgid:	426
Ngid:	0
Pid:	116220
PPid:	1
TracerPid:	0
Uid:	1000	1000	1000	1000
Gid:	1000	1000	1000	1000
FDSize:	64
Groups:
NStgid:	426
NSpid:	116220
NSpgid:	426
NSsid:	426
VmPeak:	 2167408 kB
VmSize:	 1987136 kB
VmLck:	       0 kB
VmPin:	       0 kB
VmHWM:	   62524 kB
VmRSS:	   60748 kB
VmData:	 1970764 kB
VmStk:	     132 kB
VmExe:	      76 kB
VmLib:	    3744 kB
VmPTE:	     296 kB
VmPMD:	      20 kB
VmSwap:	       0 kB
HugetlbPages:	       0 kB
Threads:	12
SigQ:	0/15593
SigPnd:	0000000000000000
ShdPnd:	0000000000000000
SigBlk:	0000000000004007
SigIgn:	0000000000001000
SigCgt:	0000000180004003
CapInh:	0000000000000000
CapPrm:	0000003fffffffff
CapEff:	0000003fffffffff
CapBnd:	0000003fffffffff
CapAmb:	0000000000000000
Seccomp:	0
Cpus_allowed:	ffffffff,ffffffff,ffffffff,ffffffff
Cpus_allowed_list:	0-127
Mems_allowed:	00000000,00000001
Mems_allowed_list:	0
voluntary_ctxt_switches:	89
nonvoluntary_ctxt_switches:	1`,
		TGID: 426,
		PID:  116220,
		Name: "vmhgfs-fuse",
		UID:  []int{1000, 1000, 1000, 1000},
		GID:  []int{1000, 1000, 1000, 1000},
	},
}

// TestStatusParse tests the parsing of /proc/PID/status files.
func TestStatusParse(t *testing.T) {
	for n, x := range statusTests {
		var s struct {
			Name string `Name`
			PID  int    `Pid`
			TGID int    `Tgid`
			UID  []int  `Uid`
			GID  []int  `Gid`
		}
		err := parseProcessStatus(bytes.NewReader([]byte(x.statusFile)),
			x.TGID, x.PID, &s)
		fmt.Printf("parsed data: %+v\n", s)
		if err != nil {
			t.Error(err)
		}
		if s.TGID != x.TGID {
			t.Errorf("TGIDs do not match in test %d (%d vs. %d)", n, s.TGID, x.TGID)
		}
		if s.PID != x.PID {
			t.Errorf("PIDs do not match in test %d (%d vs. %d)", n, s.PID, x.PID)
		}
		if s.Name != x.Name {
			t.Errorf("Names do not match in test %d (%q vs. %q)", n, s.Name, x.Name)
		}
		if len(s.UID) != len(x.UID) {
			t.Errorf("UID length mismatch in test %d (%d vs. %d)", n, len(s.UID), len(x.UID))
		} else {
			for i := range x.UID {
				if s.UID[i] != x.UID[i] {
					t.Errorf("UID[%d] mismatch in test %d (%d vs. %d)",
						i, n, s.UID[i], x.UID[i])
				}
			}
		}
		if len(s.GID) != len(x.GID) {
			t.Errorf("GID length mismatch in test %d (%d vs. %d)", n, len(s.GID), len(x.GID))
		} else {
			for i := range x.GID {
				if s.GID[i] != x.GID[i] {
					t.Errorf("GID[%d] mismatch in test %d (%d vs. %d)",
						i, n, s.GID[i], x.GID[i])
				}
			}
		}
	}
}

var cgroupTests = []struct {
	cgroupFile  string
	containerID string
}{
	{
		`13:pids:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
12:hugetlb:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
11:net_prio:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
10:perf_event:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
9:net_cls:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
8:freezer:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
7:devices:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
6:memory:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
5:blkio:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
4:cpuacct:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
3:cpu:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
2:cpuset:/docker/e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4
1:name=openrc:/docker
0::/docker
`, "e871ee9a818bab3222c94efe196e8555cb372676e96fea847a609c2d39e187a4"},
	{`10:hugetlb:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
9:perf_event:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
8:blkio:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
7:net_cls:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
6:freezer:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
5:devices:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
4:memory:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
3:cpuacct,cpu:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
2:cpuset:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
1:name=systemd:/system.slice/docker-47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81.scope
`, "47490dda5cd7e409e7bf04a8b291f87f15031090a955dac9ceed6a2160474d81"},
	{`11:hugetlb:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
10:devices:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
9:freezer:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
8:cpuacct,cpu:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
7:pids:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
6:net_prio,net_cls:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
5:blkio:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
4:perf_event:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
3:memory:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
2:cpuset:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
1:name=systemd:/kubepods/besteffort/poddbcfa688-dad5-11e7-a0e9-02e725baeeac/22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622
`, "22d8b77a1a9a6217710e3f2808c69263c674f31aa615484f808831203111e622"},
	{`9:net_cls:/
8:devices:/user.slice
7:cpu,cpuacct:/user.slice
6:pids:/user.slice/user-1000.slice/session-5.scope
5:memory:/user.slice
4:cpuset:/
3:blkio:/user.slice
2:freezer:/
1:name=systemd:/user.slice/user-1000.slice/session-5.scope
0::/user.slice/user-1000.slice/session-5.scope
`, ""},
}

func TestCgroupParse(t *testing.T) {
	for _, tc := range cgroupTests {
		cgroups := parseProcPidCgroup([]byte(tc.cgroupFile))
		cID := containerIDFromCgroups(cgroups)
		if cID != tc.containerID {
			t.Errorf("Expected container ID %s, got %s",
				tc.containerID, cID)
		}
	}
}
