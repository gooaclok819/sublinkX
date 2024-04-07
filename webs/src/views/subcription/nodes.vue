<script setup lang='ts'>
import { ref,onMounted,nextTick  } from 'vue'
import {getNodes,AddNodes,DelNode} from "@/api/subcription/node"
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
}
const tableData = ref<Node[]>([])
const Nodelink = ref('')
const Nodename = ref('')
const dialogVisible = ref(false)
const Delvisible = ref(false)
const table = ref()

onMounted(async() => {
    const {data} = await getNodes();
    tableData.value = data
})
// 格式化时间
function formatDateTime(date: Date): string {
  const y = date.getFullYear();
  const m = (date.getMonth() + 1).toString().padStart(2, '0');
  const d = date.getDate().toString().padStart(2, '0');
  const h = date.getHours().toString().padStart(2, '0');
  const min = date.getMinutes().toString().padStart(2, '0');
  const s = date.getSeconds().toString().padStart(2, '0');

  return `${y}-${m}-${d} ${h}:${min}:${s}`;
}

const addnodes = async ()=>{
   await AddNodes({
      link: Nodelink.value.trim(),
      name: Nodename.value.trim(),
    })
    tableData.value.push({
      ID: tableData.value.length + 1,
      Name: Nodename.value.trim(),
      Link: Nodelink.value.trim(),
      CreateDate: formatDateTime(new Date())
    })
    ElMessage.success("添加成功");
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

const toggleSelection = () => {
  table.value.clearSelection()
}

const handleDel = (row:any) => {
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
      tableData.value = tableData.value.filter((item) => item.ID !== row.ID)
      
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
  ).then( () => {
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
  return tableData.value.slice(start, end);
});


</script>

<template>
  <div>
    <el-dialog
    v-model="dialogVisible"
    title="添加节点"
  >
  <el-input v-model="Nodelink" placeholder="请输入节点" />
  <el-input v-model="Nodename" placeholder="请输入名称"  />
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addnodes">确定</el-button>
      </div>
    </template>
  </el-dialog>
    <el-card>
    <el-button type="primary" @click="dialogVisible=true">添加节点</el-button>
    <div style="margin-bottom: 10px"></div>
      <el-table ref="table" :data="currentTableData" style="width: 100%" @selection-change="handleSelectionChange">
    <el-table-column type="selection" fixed prop="ID" label="id"  />
    <el-table-column prop="Name" label="备注"  />
    <el-table-column prop="Link" label="节点"  />
    <el-table-column prop="CreateDate" label="创建时间" sortable  />
    <el-table-column fixed="right" label="操作" width="120">
      <template #default="scope">
        <el-button link type="primary" size="small">编辑</el-button>
  <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>

      </template>
    </el-table-column>
  </el-table>
  <div style="margin-top: 20px" />
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
.el-card{
  margin: 10px;
}
.el-input{
  margin-bottom: 10px;
}
</style>