package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"os/signal"
	"time"
	"strings"
	"github.com/prometheus/procfs"

	bpf "github.com/iovisor/gobpf/bcc"
)

const source string = `
#include <uapi/linux/ptrace.h>
#include <linux/sched.h>

struct readline_event_t {
        u32 pid;
	u32 uid;
	u64 ppid; 
        char str[256];
} __attribute__((packed));

BPF_PERF_OUTPUT(readline_events);

int get_return_value(struct pt_regs *ctx) {
	struct task_struct *task;
        struct readline_event_t event = {};
        u32 pid;
	u32 uid;
        if (!PT_REGS_RC(ctx))
                return 0;
        pid = bpf_get_current_pid_tgid();
	uid = bpf_get_current_uid_gid();
        event.pid = pid;
	event.uid = uid;
	task = (struct task_struct *)bpf_get_current_task();
	event.ppid = task->real_parent->tgid;
	bpf_probe_read(&event.str, sizeof(event.str),(void *)PT_REGS_RC(ctx));
        readline_events.perf_submit(ctx, &event, sizeof(event));

        return 0;
}
`

type readlineEvent struct {
	Pid  uint32
	Uid  uint32
	PPid  uint64
	Str  [256]byte
}

func main() {
	m := bpf.NewModule(source, []string{})
	defer m.Close()

	readlineUretprobe, err := m.LoadUprobe("get_return_value")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load get_return_value: %s\n", err)
		os.Exit(1)
	}

	err = m.AttachUretprobe("/bin/bash", "readline", readlineUretprobe, -1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to attach return_value: %s\n", err)
		os.Exit(1)
	}

	table := bpf.NewTable(m.TableId("readline_events"), m)

	channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init perf map: %s\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fmt.Printf("%20s\t%10s\t%10s\t%10s\t%10s\t%10s\t%10s\t%10s\t%10s\t%10s\n","HOATNAME","TTY","CLIENT","PID","PPID","UID","USERNAME","PWD","COMMAND","TS")
	go func() {
		var event readlineEvent
		for {
			data := <-channel
			err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}
			// Convert C string (null-terminated) to Go string
			username:=""
			cwd:=""
			hostname:=""
			tty:=""
			client:=""
			comm := string(event.Str[:bytes.IndexByte(event.Str[:], 0)])
			procObj,_:=procfs.NewProc(int(event.Pid))
			envs,_:=procObj.Environ()

			kvs:=map[string]string{}
			for _,env:=range envs{
			      kv:=strings.Split(env, "=")
			      if len(kv) == 2 {
				      kvs[kv[0]]=kv[1]
			      }
			}

			if u,ok := kvs["USER"]; ok {
			    username=u
			}else{
			    userObj,_:=user.LookupId(fmt.Sprintf("%d",event.Uid))
			    username=userObj.Username
			}

			if c,ok := kvs["PWD"]; ok {
			    cwd=c
			}else{
			    cwd,_=procObj.Cwd()
			}


			if h,ok := kvs["HOSTNAME"]; ok {
			    hostname=h
			}else{
			    hostname,_=os.Hostname()
			}

			if t,ok := kvs["SSH_TTY"]; ok {
			    tty=t
			}

			if cli,ok := kvs["SSH_CLIENT"]; ok {
			    cliObj:=strings.Split(cli, " ")
			    if len(cliObj) > 0{
				    client=cliObj[0]
			    }
			}

			t:=time.Now().Unix()
			fmt.Printf("%20s\t%10s\t%10s\t%10d\t%10d\t%10d\t%10s\t%10s\t%10s\t%10d\n",hostname,tty,client,event.Pid,event.PPid,event.Uid,username,cwd,comm,t)
		}
	}()

	perfMap.Start()
	<-sig
	perfMap.Stop()
}
