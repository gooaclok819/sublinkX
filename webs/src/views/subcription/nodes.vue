<script setup lang='ts'>
import {computed, nextTick, onMounted, ref} from 'vue'
import {AddNodes, DelNode, getNodes, SortNode, UpdateNode} from "@/api/subscription/node"
import {ElMessage, ElMessageBox} from 'element-plus'

interface Node {
  ID: number;
  Name: string;
  Link: string;
  Sort: number;
  CreateDate: string;
}

const tableData = ref<Node[]>([])
const NodeTitle = ref('')
const NodeId = ref('')
const NodeLink = ref('')
const NodeName = ref('')
const dialogVisible = ref(false)
const table = ref()
const radio1 = ref('1')

const sortDialogVisible = ref(false)
const sortedNodes = ref<Node[]>([])

async function getNodes() {
  const {data} = await getNodes();
  tableData.value = data
}

onMounted(async () => {
  await getNodes()
})

const handleAddNode = () => {
  NodeTitle.value = '添加节点'
  NodeId.value = ''
  NodeLink.value = ''
  NodeName.value = ''
  radio1.value = '1'
  dialogVisible.value = true
}

const addNodes = async () => {
  let nodeLinks = NodeLink.value.split(/[\r\n,]/);
  // 分开 过滤空行和空格
  nodeLinks = nodeLinks.map((item) => item.trim()).filter((item) => item !== '');
  if (NodeTitle.value === '添加节点') {
    // 判断合并还是分开
    if (radio1.value === '1') {
      if (NodeName.value.trim() === '') {
        ElMessage.error('备注不能为空')
        return
      }
      if (nodeLinks) {
        NodeLink.value = nodeLinks.join(',');
        await AddNodes({
          link: NodeLink.value.trim(),
          name: NodeName.value.trim(),
        })
      }
    } else {
      for (let i = 0; i < nodeLinks.length; i++) {
        await AddNodes({
          link: nodeLinks[i],
          name: "",
        })
      }
    }
    ElMessage.success("添加成功");
  } else {
    await UpdateNode({
      id: NodeId.value,
      link: NodeLink.value.trim(),
      name: NodeName.value.trim(),
    })
    ElMessage.success("更新成功");
  }
  await getNodes()
  NodeLink.value = ''
  NodeName.value = ''
  dialogVisible.value = false;
}

const multipleSelection = ref<Node[]>([])
const handleSelectionChange = (val: Node[]) => {
  multipleSelection.value = val
}

const selectAll = () => {
  nextTick(() => {
    tableData.value.forEach(row => {
      table.value.toggleRowSelection(row, true)
    })
  })
}

const handleEdit = (row: any) => {
  radio1.value = '1'
  NodeTitle.value = '编辑节点'
  NodeId.value = row.ID
  NodeName.value = row.Name
  NodeLink.value = row.Link
  dialogVisible.value = true
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
    await getNodes()
  })
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
      tableData.value = tableData.value.filter((item) => item.ID !== multipleSelection.value[i].ID)
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
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
}

// 表格数据静态化
const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableData.value.slice(start, end);
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
  sortedNodes.value = tableData.value
}
// 关闭排序弹窗时的处理
const handleCloseSortDialog = () => {
  sortDialogVisible.value = false
}
// 保存排序
const saveSortOrder = async () => {
  const response = await SortNode({
    body: JSON.stringify(tableData.value)
  })
  if (response.ok) {
    ElMessage.success('排序保存成功')
  } else {
    ElMessage.error('排序保存失败')
  }
}

</script>

<template>
  <div>
    <el-dialog
      v-model="dialogVisible"
      :title="NodeTitle"
      width="80%"
    >
      <el-input
        v-model="NodeLink"
        placeholder="请输入节点多行使用回车或逗号分开,支持base64格式的url订阅"
        type="textarea"
        style="margin-bottom:10px"
        :autosize="{ minRows: 2, maxRows: 10}"
      />
      <el-radio-group v-model="radio1" class="ml-4" v-if="NodeTitle == '添加节点'">
        <el-radio value="1" size="large">合并</el-radio>
        <el-radio value="2" size="large">分开</el-radio>
      </el-radio-group>
      <el-input v-model="NodeName" placeholder="请输入备注" v-if="radio1 != '2'"/>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="addNodes">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-card>
      <el-button type="primary" @click="handleAddNode">添加节点</el-button>
      <!-- 弹出排序弹窗的按钮 -->
      <el-button type="primary" @click="handleOpenSortDialog">修改排序</el-button>

      <!-- 排序弹窗 -->
      <el-dialog
        v-model="sortDialogVisible"
        title="调整排序"
        width="50%"
        :before-close="handleCloseSortDialog"
      >
        <draggable v-model="sortedNodes" animation="200">
          <div v-for="node in sortedNodes" :key="node.ID" class="sortable-item">
            {{ node.Name }}
          </div>
        </draggable>
        <template #footer>
          <el-button @click="sortDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveSortOrder">保存排序</el-button>
        </template>
      </el-dialog>

      <div style="margin-bottom: 10px"></div>
      <el-table ref="table" :data="currentTableData" style="width: 100%" @selection-change="handleSelectionChange">
        <el-table-column type="selection" fixed prop="ID" label="id"/>
        <el-table-column prop="Name" label="备注">
          <template #default="scope">
            <el-tag type="success">{{ scope.row.Name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Link" label="节点" sortable :show-overflow-tooltip="true"/>
        <el-table-column prop="CreateDate" label="创建时间" sortable/>
        <el-table-column fixed="right" label="操作" width="120">
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
        :total="tableData.length">
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

.sortable-item {
  padding: 10px;
  margin-bottom: 5px;
  background: #f5f5f5;
  border-radius: 4px;
  cursor: move;
}
</style>
