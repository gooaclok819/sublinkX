<script setup lang='ts'>
import { ref,onMounted,nextTick  } from 'vue'
import {getSubs,AddSub,DelSub} from "@/api/subcription/subs"
import {getNodes} from "@/api/subcription/node"
import { log } from 'console';
interface Sub {
  ID: number;
  Name: string;
  CreateDate: string;
}
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
}
interface Config {
  clash: string;
  udp: string;
  cert: string;
}
const tableData = ref<Sub[]>([])
const Clash = ref('./template/clash.yaml')
const Subname = ref('')
const dialogVisible = ref(false)
const Delvisible = ref(false)
const table = ref()
const NodesList = ref<Node[]>([])
const value1 = ref([])
const checkList = ref<string[]>([]) // 配置列表

onMounted(async() => {
    const {data} = await getSubs();
    tableData.value = data
})
onMounted(async() => {
    const {data} = await getNodes();
    NodesList.value = data
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


const addSubs = async ()=>{
    const config = JSON.stringify({
    "clash": Clash.value.trim(),
    "udp": checkList.value.includes('udp') ? true :  false,
    "cert": checkList.value.includes('cert') ? true :  false

  })
   await AddSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(',')
    })
    tableData.value.push({
      ID: tableData.value.length + 1,
      Name: Subname.value.trim(),
      CreateDate: formatDateTime(new Date())
    })
    ElMessage.success("添加成功");
    dialogVisible.value = false;
}

const multipleSelection = ref<Sub[]>([])
const handleSelectionChange = (val: Sub[]) => {
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
      await DelSub({
        id: row.ID
      })
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      tableData.value = tableData.value.filter((item) => item.ID !== row.ID)
      
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
  ).then( () => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
       DelSub({
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
}
// 表格数据静态化
const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return tableData.value.slice(start, end);
});
const isNameShow = (row: any): boolean =>  {
  return row.Name.length > 10;
}

</script>

<template>
  <div>
    <el-dialog
    v-model="dialogVisible"
    title="添加订阅"
  >
  <el-input v-model="Subname" placeholder="请输入订阅名称" />
  
  <el-row >
  <el-tag type="primary">clash模版文件</el-tag>
  <el-input v-model="Clash" placeholder="clash模版文件"  />
</el-row>

  <el-row>
    <el-tag type="primary">强制开启选项</el-tag>
  <el-checkbox-group v-model="checkList"  style="margin: 5px;">
    <el-checkbox :value="'udp'">udp</el-checkbox>
    <el-checkbox :value="'cert'">跳过证书</el-checkbox>
  </el-checkbox-group>
</el-row>
  <div class="m-4">
    <p>选择已有的节点列表</p>
    <el-select
      v-model="value1"
      multiple
      placeholder="Select"
      style="width: 100%"
    >
      <el-option
        v-for="item in NodesList"
        :key="item.Name"
        :label="item.Name"
        :value="item.Name"
      />
    </el-select>
  </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="addSubs">确定</el-button>
      </div>
    </template>
  </el-dialog>
    <el-card>
    <el-button type="primary" @click="dialogVisible=true">添加节点</el-button>
    <div style="margin-bottom: 10px"></div>

      <el-table ref="table" 
      :data="currentTableData" 
      style="width: 100%" 
      stripe
      @selection-change="handleSelectionChange" 
      row-key="id" 
      :tree-props="{children: 'Nodes'}"
      >
    <el-table-column type="selection" fixed prop="ID" label="id"  />
    <el-table-column prop="Name" label="订阅名称"  >
    <template #default="{row}">
      <el-tag :type="row.Nodes ? 'success' : 'primary'">{{row.Name}}</el-tag>
        </template>
    </el-table-column>
    <el-table-column prop="Link" label="链接" :show-overflow-tooltip="true" />
 
    <el-table-column prop="CreateDate" label="创建时间" sortable  />
    <el-table-column  label="操作" width="120">
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
.el-tag{
  margin: 5px;
}

</style>