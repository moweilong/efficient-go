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

// TestXORFeatures 测试异或运算的两个核心特性：翻转特定位、判断符号是否相同
func TestXORFeatures(t *testing.T) {
	// 测试1：异或运算翻转特定位（前8位，从MSB开始）
	t.Run("flip_bits", func(t *testing.T) {
		var a uint16 = 0xCEFF  // 二进制：11001110 11111111
		mask := uint16(0xFF00) // 掩码：11111111 00000000（前8位为1）
		a ^= mask              // 异或运算：翻转前8位

		expected := uint16(0x31FF) // 预期结果：00110001 11111111
		if a != expected {
			t.Errorf("特定位翻转失败\n原始值: 0x%X\n掩码: 0x%X\n实际结果: 0x%X\n预期结果: 0x%X",
				0xCEFF, mask, a, expected)
		}
		t.Logf("特定位翻转成功\n原始值: 0x%X → 异或0x%X → 结果: 0x%X", 0xCEFF, mask, a)
	})

	// 测试2：异或运算判断两个整数的符号是否相同
	t.Run("check_sign", func(t *testing.T) {
		// 测试用例：(a, b, 预期符号是否相同)
		testCases := []struct {
			name     string
			a, b     int
			expected bool
		}{
			{"均为正数", 12, 25, true},
			{"均为负数", -12, -25, true},
			{"一正一负", -12, 25, false},
			{"一负一正", 12, -25, false},
			{"包含零（零视为正数）", 0, -5, false},
			{"零与正数", 0, 5, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// 核心逻辑：(a ^ b) >= 0 表示符号相同
				actual := (tc.a ^ tc.b) >= 0
				if actual != tc.expected {
					t.Errorf("符号判断错误\na=%d, b=%d\n实际: %v\n预期: %v",
						tc.a, tc.b, actual, tc.expected)
				}
				t.Logf("符号判断正确\na=%d, b=%d → 符号%s",
					tc.a, tc.b, map[bool]string{true: "相同", false: "不同"}[actual])
			})
		}
	})
}

// TestBitwiseNot 测试Go中的一元按位取反运算符^
// 验证^a对变量a的所有位进行取反操作（0变1，1变0）
func TestBitwiseNot(t *testing.T) {
	// 测试用例：原始值、预期取反结果（针对byte类型）
	testCases := []struct {
		name     string
		input    byte
		expected byte
	}{
		{
			name:     "0x0F取反",
			input:    0x0F, // 二进制：00001111
			expected: 0xF0, // 二进制：11110000
		},
		{
			name:     "0x00取反",
			input:    0x00, // 二进制：00000000
			expected: 0xFF, // 二进制：11111111
		},
		{
			name:     "0xFF取反",
			input:    0xFF, // 二进制：11111111
			expected: 0x00, // 二进制：00000000
		},
		{
			name:     "0xAB取反",
			input:    0xAB, // 二进制：10101011
			expected: 0x54, // 二进制：01010100
		},
	}

	// 测试一元^的取反功能
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ^tc.input // 执行按位取反操作

			// 验证结果是否符合预期
			if result != tc.expected {
				t.Errorf("取反结果错误\n原始值: 0x%X (%08b)\n实际结果: 0x%X (%08b)\n预期结果: 0x%X (%08b)",
					tc.input, tc.input,
					result, result,
					tc.expected, tc.expected)
			} else {
				t.Logf("取反成功\n原始值: 0x%X (%08b) → 取反后: 0x%X (%08b)",
					tc.input, tc.input,
					result, result)
			}
		})
	}

	// 单独测试单个位的翻转（利用二元^运算符）
	t.Run("单个位翻转", func(t *testing.T) {
		// 测试用例：(原始位值, 翻转后预期值)
		bitTests := []struct {
			bit      int
			expected int
		}{
			{0, 1}, // 0与1异或 → 1
			{1, 0}, // 1与1异或 → 0
		}

		for _, bt := range bitTests {
			result := 1 ^ bt.bit // 用二元^实现单个位翻转
			if result != bt.expected {
				t.Errorf("位翻转错误\n原始位: %d → 翻转后: %d, 预期: %d",
					bt.bit, result, bt.expected)
			} else {
				t.Logf("位翻转成功\n原始位: %d → 翻转后: %d", bt.bit, result)
			}
		}
	})
}

// TestAndNotOperator 测试按位与非运算符&^的功能
// 验证&^是否能根据第二个操作数清除第一个操作数的特定位
func TestAndNotOperator(t *testing.T) {
	// 测试用例：(原始值a, 掩码b, 预期结果)
	testCases := []struct {
		name     string
		a        byte
		b        byte
		expected byte
	}{
		{
			name:     "清除低4位",
			a:        0xAB, // 二进制：10101011
			b:        0x0F, // 掩码：00001111（低4位为1，用于清除）
			expected: 0xA0, // 结果：10100000（低4位被清除为0）
		},
		{
			name:     "清除高4位",
			a:        0xAB, // 10101011
			b:        0xF0, // 11110000（高4位为1）
			expected: 0x0B, // 00001011（高4位被清除为0）
		},
		{
			name:     "不清除任何位（b全为0）",
			a:        0xAB, // 10101011
			b:        0x00, // 00000000
			expected: 0xAB, // 保留原始值
		},
		{
			name:     "清除所有位（b全为1）",
			a:        0xAB, // 10101011
			b:        0xFF, // 11111111
			expected: 0x00, // 所有位被清除为0
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 执行&^=操作（等价于 a = a &^ tc.b）
			result := tc.a
			result &^= tc.b

			// 验证结果
			if result != tc.expected {
				t.Errorf("与非运算结果错误\n原始值a: 0x%X (%08b)\n掩码b: 0x%X (%08b)\n实际结果: 0x%X (%08b)\n预期结果: 0x%X (%08b)",
					tc.a, tc.a,
					tc.b, tc.b,
					result, result,
					tc.expected, tc.expected)
			} else {
				t.Logf("与非运算成功\n原始值a: 0x%X (%08b) &^ 掩码b: 0x%X (%08b) → 结果: 0x%X (%08b)",
					tc.a, tc.a,
					tc.b, tc.b,
					result, result)
			}
		})
	}
}

// TestShiftOperators 测试移位运算符<<和>>的功能及应用
func TestShiftOperators(t *testing.T) {
	// 子测试1：左移基本功能（无符号数）
	t.Run("left_shift_basic", func(t *testing.T) {
		var a int8 = 3 // 二进制：00000011
		testCases := []struct {
			shift  int
			expect int8
			binary string
		}{
			{1, 6, "00000110"},
			{2, 12, "00001100"},
			{3, 24, "00011000"},
		}
		for _, tc := range testCases {
			result := a << tc.shift
			if result != tc.expect {
				t.Errorf("左移%d位错误\n原始值: %d (%08b)\n实际结果: %d (%08b)\n预期结果: %d (%s)",
					tc.shift, a, a, result, result, tc.expect, tc.binary)
			}
		}
	})

	// 子测试2：右移基本功能（无符号数，逻辑移位）
	t.Run("right_shift_unsigned", func(t *testing.T) {
		var a uint8 = 120 // 二进制：01111000
		testCases := []struct {
			shift  int
			expect uint8
			binary string
		}{
			{1, 60, "00111100"},
			{2, 30, "00011110"},
		}
		for _, tc := range testCases {
			result := a >> tc.shift
			if result != tc.expect {
				t.Errorf("右移%d位错误（无符号数）\n原始值: %d (%08b)\n实际结果: %d (%08b)\n预期结果: %d (%s)",
					tc.shift, a, a, result, result, tc.expect, tc.binary)
			}
		}
	})

	// 子测试3：右移（有符号数，算术移位）
	t.Run("right_shift_signed", func(t *testing.T) {
		var a int8 = -8 // 二进制：11111000（补码）
		testCases := []struct {
			shift  int
			expect int8
			binary string // 算术移位：高位补符号位1
		}{
			{1, -4, "11111100"},
			{2, -2, "11111110"},
		}
		for _, tc := range testCases {
			result := a >> tc.shift
			if result != tc.expect {
				t.Errorf("右移%d位错误（有符号数）\n原始值: %d (%08b)\n实际结果: %d (%08b)\n预期结果: %d (%s)",
					tc.shift, a, a, result, result, tc.expect, tc.binary)
			}
		}
	})

	// 子测试4：移位与乘除法（2^n）
	t.Run("shift_mult_div", func(t *testing.T) {
		// 左移 = 乘法（a << n = a * 2^n）
		{
			a := 12
			result := a << 2 // 12 * 4 = 48
			if result != 48 {
				t.Errorf("左移2位乘法错误\n原始值: %d, 实际结果: %d, 预期: 48", a, result)
			}
		}
		// 右移 = 除法（a >> n = a / 2^n）
		{
			a := 200
			result := a >> 1 // 200 / 2 = 100
			if result != 100 {
				t.Errorf("右移1位除法错误\n原始值: %d, 实际结果: %d, 预期: 100", a, result)
			}
		}
	})

	// 子测试5：结合|和<<设置特定位（第n位）
	t.Run("set_bit_with_or_shift", func(t *testing.T) {
		var a int8 = 8 // 二进制：00001000
		n := 2         // 要设置的位（从0开始，第2位）
		a |= 1 << n    // 00001000 | 00000100 = 00001100
		expected := int8(12)
		if a != expected {
			t.Errorf("设置第%d位错误\n结果: %d (%08b), 预期: %d (00001100)",
				n, a, a, expected)
		}
	})

	// 子测试6：结合&和<<测试特定位是否设置
	t.Run("check_bit_with_and_shift", func(t *testing.T) {
		var a int8 = 12 // 二进制：00001100（第2位为1）
		n := 2
		isSet := a&(1<<n) != 0 // 00001100 & 00000100 = 00000100 ≠ 0 → true
		if !isSet {
			t.Errorf("测试第%d位错误\n预期: 已设置, 实际: 未设置", n)
		}
	})

	// 子测试7：结合&^和<<清除特定位
	t.Run("clear_bit_with_andnot_shift", func(t *testing.T) {
		var a int8 = 13 // 二进制：00001101（第2位为1）
		n := 2
		a &^= 1 << n // 00001101 &^ 00000100 = 00001001
		expected := int8(9)
		if a != expected {
			t.Errorf("清除第%d位错误\n结果: %d (%08b), 预期: %d (00001001)",
				n, a, a, expected)
		}
	})
}
