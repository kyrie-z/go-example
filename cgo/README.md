
## 类型转换

### GO references to C
1. 指针互转
  - 通过**unsafe.Pointer**类型实现GO和C的指针，类似`void*`.任何C指针都可转为unsafe.Pointer。
  - 若要进行指针运算，如结构体偏移，需要先用unsafe.Pointer后再使用uintptr后对指针地址进行运算。如修改结构体的第二个与元素(x.b)int型，`pb := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))`，pb为第二个元素的地址，再进行修改。**不要试图引入一个uintptr类型的临时变量**。GC优化会导致变量地址偏移。

2. 字符串(string)与字节切片([]byte)
  - GO string 转 `char*` 如：`cs := C.CString(gstring)`，需要手动释放`defer C.free(unsafe.Pointer(cs))`
  - GO []byte 转 unsafe.Pointer 再转C指针 `char*` 如： `var cb = C.CBytes(gbyte)`，`(*C.char)(cb)`，最后需要释放`defer C.free(cb)`
  - 当需要传入C `char*`类型数据用于修改时，直接传递[]byte给C，`unsafe.Pointer(&bytes[0])`

3. 数值和指针的转换
  - Go语言有uintptr类型。以uintptr为中介，实现数值类型到unsafe.Pointr指针类型到转换，如int转`void*` ` unsafe.Pointer(uintptr(hKey)`，在通过`void*`转其他C指针类型。

4. 结构体、联合、枚举类型(C中未使用typedef)
  - **通过C.struct_xxx来访问C语言中定义的struct xxx结构体类型**,如`var h *C.struct_crypt_cipher`;
  - 如果结构体的成员名字中碰巧是`Go语言的关键字`，通过在成员名开头添加下划线来访问,如`a._type`.

5. 访问C内存空间
  - 访问C数组`char arr[10]` ，如`arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]`
  - 访问C结构体 ，如`rBuf := (*[1 << 32]byte)(unsafe.Pointer(&cSD.r[0]))[:int(32):int(32)]`
  - 访问C字符串`char *s = "Hello";`，如`s1 := string((*[31]byte)(unsafe.Pointer(C.s))[:sLen:sLen])` sLen为字符串长度`sLen := int(C.strlen(C.s))`
 


## Function
// Go字符串转换为C字符串。C字符串使用malloc分配，因此需要使用C.free以避免内存泄露 \
func C.CString(string) *C.char

// Go byte数组转换为C的数组。使用malloc分配的空间，因此需要使用C.free避免内存泄漏 \
func C.CBytes([]byte) unsafe.Pointer

// C字符串转换为Go字符串 \
func C.GoString(*C.char) string

// C字符串转换为Go字符串，指定转换长度 \
func C.GoStringN(*C.char, C.int) string

// C数据转换为byte数组，指定转换的长度 \ 
func C.GoBytes(unsafe.Pointer, C.int) []byte
