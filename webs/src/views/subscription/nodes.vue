<script setup lang='ts'>
import {nextTick, onMounted, ref} from 'vue'
import {AddNodes, DelNode, getNodes, SortNode, UpdateNode} from "@/api/subscription/node"

interface Node {
  ID: number;
  Name: string;
  Link: string;
  Sort: number;
  CreateDate: string;
}

// 表格节点数据
const tableNodes = ref<Node[]>([])
// 排序节点数据
const sortedNodes = ref<Node[]>([])
// 编辑窗口的标题
const editNodeTitle = ref('')
// 编辑窗口的节点ID
const editNodeId = ref('')
// 输入的节点链接
const inputNodeLink = ref('')
// 输入的节点名称
const inputNodeName = ref('')
// 节点是否合并
const nodeMergingRadio = ref('1')
// 编辑窗口
const editDialogVisible = ref(false)
// 排序窗口
const sortDialogVisible = ref(false)
// 排序树显示的名称
const sortTreeProps = {
  // 指定显示的名称字段
  label: 'Name',
};
const table = ref()

async function fetchNodes() {
  const {data} = await getNodes();
  tableNodes.value = data
}

onMounted(async () => {
  await fetchNodes()
})

const handleAddNode = () => {
  editNodeTitle.value = '添加节点'
  inputNodeLink.value = ''
  inputNodeName.value = ''
  nodeMergingRadio.value = '1'
  editDialogVisible.value = true
}

const editNodes = async () => {
  let splitNodeLinks = inputNodeLink.value.split(/[\r\n,]/);
  // 分开 过滤空行和空格
  splitNodeLinks = splitNodeLinks.map((item) => item.trim()).filter((item) => item !== '');
  if (editNodeTitle.value == '添加节点') {
    // 判断合并还是分开
    if (nodeMergingRadio.value === '1') {
      if (inputNodeName.value.trim() === '') {
        ElMessage.error('备注不能为空')
        return
      }
      if (splitNodeLinks) {
        inputNodeLink.value = splitNodeLinks.join(',');
        await AddNodes({
          link: inputNodeLink.value.trim(),
          name: inputNodeName.value.trim(),
        })
      }
    } else {
      for (let i = 0; i < splitNodeLinks.length; i++) {
        await AddNodes({
          link: splitNodeLinks[i],
          name: "",
        })
      }
    }
    ElMessage.success("添加成功");
  } else {
    await UpdateNode({
      id: editNodeId.value,
      link: inputNodeLink.value.trim(),
      name: inputNodeName.value.trim(),
    })
    ElMessage.success("更新成功");
  }
  await fetchNodes()
  inputNodeLink.value = ''
  inputNodeName.value = ''
  editDialogVisible.value = false;
}

const multipleSelection = ref<Node[]>([])
const handleSelectionChange = (val: Node[]) => {
  multipleSelection.value = val

}
const selectAll = () => {
  nextTick(() => {
    tableNodes.value.forEach(row => {
      table.value.toggleRowSelection(row, true)
    })
  })
}
const handleEdit = (row: any) => {
  nodeMergingRadio.value = '1'
  editNodeTitle.value = '编辑节点'
  editNodeId.value = row.ID
  inputNodeName.value = row.Name
  inputNodeLink.value = row.Link
  editDialogVisible.value = true
}
const toggleSelection = () => {
  table.value.clearSelection()
}

const handleDel = (row: any) => {
  ElMessageBox.confirm(
    `你是否要删除 ${row.Name} ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(async () => {
    await DelNode({
      id: row.ID
    })
    ElMessage({
      type: 'success',
      message: '删除成功',
    })
    await fetchNodes()
    // tableData.value = tableData.value.filter((item) => item.ID !== row.ID)

  })
  // console.log('click',row.ID)
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) {
    return
  }
  ElMessageBox.confirm(
    `你是否要删除选中这些 ?`,
    '提示',
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(() => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
      DelNode({
        id: multipleSelection.value[i].ID
      })
      tableNodes.value = tableNodes.value.filter((item) => item.ID !== multipleSelection.value[i].ID)
    }
    ElMessage({
      type: 'success',
      message: '删除成功',
    })
  })

}
// 分页显示
const currentPage = ref(1);
const pageSize = ref(10);
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  // console.log(`每页 ${val} 条`);
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  // console.log(`当前页: ${val}`);
}
// 表格数据静态化
const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableNodes.value.slice(start, end);
});
// 复制链接
const copyUrl = (url: string) => {
  const textarea = document.createElement('textarea');
  textarea.value = url;
  document.body.appendChild(textarea);
  textarea.select();
  try {
    const successful = document.execCommand('copy');
    const msg = successful ? 'success' : 'warning';
    const message = successful ? '复制成功！' : '复制失败！';
    ElMessage({
      type: msg,
      message,
    });
  } catch (err) {
    ElMessage({
      type: 'warning',
      message: '复制失败！',
    });
  } finally {
    document.body.removeChild(textarea);
  }
};
const copyInfo = (row: any) => {
  copyUrl(row.Link)
}

// 打开排序弹窗时初始化排序数据
const handleOpenSortDialog = () => {
  sortDialogVisible.value = true
  getNodes()
    .then(({data}) => {
      sortedNodes.value = data;
    })
    .catch(error => {
      console.error('Failed to fetch nodes:', error);
    });
}
// 关闭排序弹窗时的处理
const handleCloseSortDialog = () => {
  sortDialogVisible.value = false
}

// 节点拖拽后的处理逻辑
const allowNodeDrop = (draggingNode: any, dropNode: any, dropType: string) => {
  return dropType !== 'inner';
}

const handleNodeDrop = (draggingNode: any, dropNode: any, dropType: string) => {
  // 确保只有允许的拖拽操作才会触发排序更新
  if (dropType === 'inner') {
    ElMessage.error('不允许的拖拽操作');
    return;
  }

  // 根据当前的顺序重新设置 Sort 字段
  sortedNodes.value.forEach((node, index) => {
    node.Sort = index + 1;  // 索引从 1 开始
  });
}
// 保存排序
const saveSortOrder = async () => {
  try {
    await SortNode(sortedNodes.value);
    // 正确地将数据放在 data 属性中
    ElMessage.success('排序保存成功');
  } catch (error) {
    ElMessage.error('请求失败');
  } finally {
    handleCloseSortDialog();
    await fetchNodes();
  }
}

</script>

<template>
  <div>
    <!-- 编辑弹窗 -->
    <el-dialog v-model="editDialogVisible" :title="editNodeTitle" width="70%">
      <el-input
        v-model="inputNodeLink"
        placeholder="请输入节点多行使用回车或逗号分开,支持base64格式的url订阅"
        type="textarea"
        style="margin-bottom:20px"
        :autosize="{ minRows: 3, maxRows: 10}"
      />
      <el-radio-group v-model="nodeMergingRadio" class="ml-4" v-if="editNodeTitle== '添加节点'">
        <el-radio value="1" size="large">合并</el-radio>
        <el-radio value="2" size="large">分开</el-radio>
      </el-radio-group>
      <el-input v-model="inputNodeName" placeholder="请输入备注" v-if="nodeMergingRadio!='2'"/>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="editNodes">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 排序弹窗 -->
    <el-dialog v-model="sortDialogVisible" title="排序节点" width="40%">
      <el-tree
        :data="sortedNodes"
        node-key="ID"
        default-expand-all
        draggable
        :props="sortTreeProps"
        :allow-drop="allowNodeDrop"
        @node-drop="handleNodeDrop"
      >
      </el-tree>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleCloseSortDialog">关闭</el-button>
          <el-button type="primary" @click="saveSortOrder">保存排序</el-button>
        </div>
      </template>
    </el-dialog>

    <el-card>
      <el-button type="primary" @click="handleAddNode">添加节点</el-button>
      <!-- 弹出排序弹窗的按钮 -->
      <el-button type="primary" @click="handleOpenSortDialog">修改排序</el-button>
      <div style="margin-bottom: 10px"></div>
      <el-table ref="table" :data="currentTableData" style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" fixed prop="ID" label="id"/>
        <el-table-column prop="Name" label="备注" sortable width="200">
          <template #default="scope">
            <el-tag type="success">{{ scope.row.Name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Link" label="节点" :show-overflow-tooltip="true"/>
        <el-table-column prop="CreateDate" label="创建时间" sortable width="200"/>
        <el-table-column fixed="right" label="操作" width="240">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="copyInfo(scope.row)">复制</el-button>
            <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 20px"/>
      <el-button type="info" @click="selectAll()">全选</el-button>
      <el-button type="warning" @click="toggleSelection()">取消选择</el-button>
      <el-button type="danger" @click="selectDel">批量删除</el-button>
      <div style="margin-top: 20px"/>
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 30, 40]"
        :total="tableNodes.length">
      </el-pagination>

    </el-card>
  </div>
</template>

<style scoped>
.el-card {
  margin: 10px;
}
.el-input {
  margin-bottom: 10px;
}
</style>
