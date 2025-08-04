package bit_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TestClearLast4LSB 使用 & 运算符清除最后 4 个最低有效位（LSB）为 0
func TestClearLast4LSB(t *testing.T) {
	var x uint8 = 0xAC // x = 10101100（二进制）
	x = x & 0xF0       // 与 11110000（0xF0）进行与运算，结果为 10100000（0xA0）

	// 验证结果：0xF0 与 0xAC 进行与运算后应得到 0xA0
	expected := uint8(0xA0)
	if x != expected {
		t.Errorf("位运算结果错误: 期望 0x%X，实际得到 0x%X", expected, x)
	}
	t.Logf("位运算结果: 期望 0x%X，实际得到 0x%X", expected, x)
}

// TestBitwiseShortAssignmentForLSB 使用 &= 简写形式清除最后 4 个最低有效位（LSB）为 0
func TestBitwiseShortAssignmentForLSB(t *testing.T) {
	var x uint8 = 0xAC // x = 10101100（二进制）
	// &= 是 & 运算符的简写赋值形式（等价于 x = x & 0xF0）
	x &= 0xF0 // 结果为 10100000（0xA0）

	// 验证结果：0xF0 与 0xAC 进行与运算后应得到 0xA0
	expected := uint8(0xA0)
	if x != expected {
		t.Errorf("位运算结果错误: 期望 0x%X，实际得到 0x%X", expected, x)
	}
	t.Logf("位运算结果: 期望 0x%X，实际得到 0x%X", expected, x)
}

// TestOddEvenWithBitwise 用 & 操作测试一个数字是奇数还是偶数
// 原理：任何整数与1进行与运算，结果为1则是奇数，为0则是偶数
func TestOddEvenWithBitwise(t *testing.T) {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 测试100个随机数
	for i := range make([]struct{}, 100) {
		num := r.Int() // 使用新生成器产生随机数
		// 使用位运算判断奇偶
		isOdd := num&1 == 1

		// 使用取模运算作为对照验证
		expectedOdd := num%2 == 1

		// 记录测试结果
		t.Logf("测试第%d个数字: %d, 位运算判断为%s",
			i+1, num, map[bool]string{true: "奇数", false: "偶数"}[isOdd])

		// 验证位运算结果是否正确
		if isOdd != expectedOdd {
			t.Errorf("判断错误: 数字%d, 位运算判断为%s, 实际为%s",
				num,
				map[bool]string{true: "奇数", false: "偶数"}[isOdd],
				map[bool]string{true: "奇数", false: "偶数"}[expectedOdd])
		}
	}
}

// TestSetBitsWithOr 用按位或（|）操作有选择地设置整数的特定位为1
// 目标：将uint8类型变量的第3位、第7位、第8位（从低位开始计数）设为1
func TestSetBitsWithOr(t *testing.T) {
	var a uint8 = 0 // 初始值：00000000（二进制）

	// 构造掩码：需要设置的位为1，其余为0
	// 第3位（2^2）、第7位（2^6）、第8位（2^7）对应的掩码计算：
	// 2^2 = 4，2^6 = 64，2^7 = 128 → 掩码 = 4 + 64 + 128 = 196
	mask := uint8(196) // 掩码二进制：11000100（第3、7、8位为1）

	a |= mask // 按位或操作：0 | 196 = 196（二进制11000100）

	// 预期结果：196（二进制11000100）
	expected := uint8(196)
	if a != expected {
		t.Errorf("按位或结果错误：期望 %b（0x%X），实际 %b（0x%X）",
			expected, expected, a, a)
	}

	// 验证特定位是否被正确设置（分别检查第3、7、8位）
	// 检查方法：用对应位的掩码与结果做&运算，若不为0则该位为1
	bitsToCheck := []struct {
		bit  int   // 位序号（从1开始）
		mask uint8 // 该位对应的掩码（2^(bit-1)）
	}{
		{bit: 3, mask: 1 << 2}, // 第3位：2^2 = 4（二进制00000100）
		{bit: 7, mask: 1 << 6}, // 第7位：2^6 = 64（二进制01000000）
		{bit: 8, mask: 1 << 7}, // 第8位：2^7 = 128（二进制10000000）
	}

	for _, check := range bitsToCheck {
		if (a & check.mask) == 0 {
			t.Errorf("第%d位未被正确设置为1：实际值为0", check.bit)
		} else {
			t.Logf("第%d位设置正确：值为1", check.bit)
		}
	}

	t.Logf("最终结果：%b（0x%X）", a, a)
}

// 定义位掩码常量（每个常量对应一个2的幂，确保二进制中只有一位为1）
const (
	UPPER = 1 << iota // 1 << 0 = 1（二进制0001）：转换为大写
	LOWER             // 1 << 1 = 2（二进制0010）：转换为小写
	CAP               // 1 << 2 = 4（二进制0100）：单词首字母大写
	REV               // 1 << 3 = 8（二进制1000）：反转字符串
)

// procstr 根据位掩码配置对字符串执行转换操作
// 参数：
//
//	str: 待转换的字符串
//	conf: 位掩码配置（通过UPPER|LOWER|CAP|REV组合）
//
// 返回：转换后的字符串
func procstr(str string, conf byte) string {
	// 反转字符串的内部函数
	rev := func(s string) string {
		runes := []rune(s)
		n := len(runes)
		for i := 0; i < n/2; i++ {
			runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
		}
		return string(runes)
	}

	// 根据位掩码配置执行对应操作（用&判断特定位是否被设置）
	if (conf & UPPER) != 0 {
		str = strings.ToUpper(str)
	}
	if (conf & LOWER) != 0 {
		str = strings.ToLower(str)
	}
	if (conf & REV) != 0 {
		str = rev(str)
	}
	if (conf & CAP) != 0 {
		str = cases.Title(language.English).String(str)
	}
	return str
}

// TestBitmaskConfig 测试位掩码技术在多配置场景中的应用
// 验证通过|组合配置项、通过&查询配置项的正确性
func TestBitmaskConfig(t *testing.T) {
	// 测试用例：输入字符串和预期结果
	testCases := []struct {
		name     string
		input    string
		conf     byte
		expected string
	}{
		{
			name:     "LOWER|REV|CAP组合",
			input:    "HELLO PEOPLE!",
			conf:     LOWER | REV | CAP, // 2|8|4=14（二进制1110）
			expected: "!Elpoep Olleh",
		},
		{
			name:     "UPPER|REV组合",
			input:    "hello",
			conf:     UPPER | REV, // 1|8=9（二进制1001）
			expected: "OLLEH",
		},
		{
			name:     "无配置",
			input:    "test",
			conf:     0,
			expected: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := procstr(tc.input, tc.conf)
			if result != tc.expected {
				t.Errorf("配置%s处理失败\n输入: %q\n预期: %q\n实际: %q",
					tc.name, tc.input, tc.expected, result)
				return
			}
			t.Logf("配置%s处理成功\n输入: %q\n配置值: %d（二进制%04b）\n结果: %q",
				tc.name, tc.input, tc.conf, tc.conf, result)
		})
	}
}
