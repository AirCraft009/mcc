#define video_char_table_start:49185
#define video_char_table_size:2048

#define video_start:32768
#define video_end:49151

#define Dimension:256        // native resolution is 256x256 (upscaled elsewhere)

#define Bpp:2
#define Ppb:4               // pixels per byte

#define stack_start:32767
#define split_stack_size:910
#define task_size:61        // real size 60, +1 for easier math

#define active_task_location:9088

#define task_start:9089
#define task_len_pos:9149

#define ReadPtr:49152
#define WritePtr:49153

#define bitmap_start:8192
#define bitmap_end:9087

#define writeable_heap:9629

_get_active_task:
    MOVI T6 9088
    LOADB T6 T6
    RET

_get_task_len:
    MOVI T6 9149
    LOADB T6 T6
    RET

_get_state_location:
    CALL _get_task_size
    MUL T1 T6       # get Offsets
    CALL _get_task_start
    ADD T1 T6       # get location
    SUBI T1 2
    RET

_get_state:          # T1 has the task number from 0
    CALL _get_state_location
    LOADB T1 T1     # LOAD into T1
    RET