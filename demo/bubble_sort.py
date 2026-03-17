# 冒泡排序算法实现

def bubble_sort(arr):
    """
    对输入的列表进行冒泡排序
    :param arr: 待排序的列表
    :return: 排序后的列表
    """
    n = len(arr)
    # 外层循环控制排序轮数
    for i in range(n):
        # 标志位，用于判断是否发生交换，若未发生交换则说明已排序完成
        swapped = False
        # 内层循环进行相邻元素比较和交换
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                # 交换元素
                arr[j], arr[j + 1] = arr[j + 1], arr[j]
                swapped = True
        # 如果没有发生交换，说明数组已经有序，可以提前结束
        if not swapped:
            break
    return arr

# 示例用法
if __name__ == "__main__":
    example_list = [64, 34, 25, 12, 22, 11, 90]
    print("原始列表:", example_list)
    sorted_list = bubble_sort(example_list.copy())
    print("排序后列表:", sorted_list)
