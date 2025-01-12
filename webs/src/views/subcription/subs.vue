<script setup lang='ts'>
import { ref, onMounted } from "vue";
import {getSubs,AddSub,DelSub,UpdateSub,UpdateSubSort} from "@/api/subcription/subs"
import {getTemp} from "@/api/subcription/temp"
import {getNodes} from "@/api/subcription/node"
import QrcodeVue from 'qrcode.vue'
import md5 from 'md5'
import draggable from 'vuedraggable';

interface Sub {
  ID: number;
  Name: string;
  CreateDate: string;
  Config: Config;
  Nodes: Node[];
  SubLogs:SubLogs[];
}
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
  sort?: number; // 添加 sort 字段
}
interface Config {
  clash: string;
  surge:string;
  udp: string;
  cert: string;
}
interface SubLogs {
  date: string;
  name: string;
  count: number;
  address: string;
}
interface Temp {
  file: string;
  text: string;
  CreateDate: string;
}
const tableData = ref<Sub[]>([])
const Clash = ref('')
const Surge = ref('')
const SubTitle = ref('')
const Subname = ref('')
const oldSubname = ref('')
const dialogVisible = ref(false)
const table = ref()
const NodesList = ref<Node[]>([])
const value1 = ref<string[]>([])
const sortSubID = ref()
const sortValue = ref<Node[]>([])
const checkList = ref<string[]>([]) // 配置列表
const iplogsdialog = ref(false)
const IplogsList = ref<SubLogs[]>([])
const qrcode = ref('')
const templist = ref<Temp[]>([])
async function getsubs() {
  const {data} = await getSubs();
    tableData.value = data
}
async function gettemps() {
    const {data} = await getTemp();
    templist.value = data
    console.log(templist.value);
}
onMounted(() => {
    getsubs()
    gettemps()
})
onMounted(async() => {
    const {data} = await getNodes();
    NodesList.value = data
})

const addSubs = async ()=>{
    const config = JSON.stringify({
    "clash": Clash.value.trim(),
    "surge": Surge.value.trim(),
    "udp": checkList.value.includes('udp') ? true :  false,
    "cert": checkList.value.includes('cert') ? true :  false

  })
  if (SubTitle.value === '添加订阅') {
    await AddSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(',')
    })
    getsubs()
    ElMessage.success("添加成功");
  }else{
    await UpdateSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(','),
      oldname: oldSubname.value
    })
    getsubs()
    ElMessage.success("更新成功");
  }

    dialogVisible.value = false;
}

const multipleSelection = ref<Sub[]>([])
const handleSelectionChange = (val: Sub[]) => {
  multipleSelection.value = val

}
const selectAll = () => {
  tableData.value.forEach(row => {
            table.value.toggleRowSelection(row, true)
        })
}
const handleIplogs = (row: any) => {
  iplogsdialog.value = true
  nextTick(() => {
    tableData.value.forEach((item) => {
    if (item.ID === row.ID) {
      IplogsList.value = item.SubLogs
    }
  })

  })
}

const toggleSelection = () => {
  table.value.clearSelection()
}

const handleAddSub = ()=>{
  SubTitle.value = '添加订阅'
  Subname.value = ''
  oldSubname.value = ''
  checkList.value = []
  Clash.value = './template/clash.yaml'
  Surge.value = './template/surge.conf'
  dialogVisible.value = true
  value1.value = []
  sortSubID.value = null; // 新增页面没有 subId
  sortValue.value = []; // 清空排序数据
}
const handleEdit = (row:any) => {
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].ID === row.ID) {
      function toConfig(value: string | Config): Config {
        if (typeof value === 'string') {
          return JSON.parse(value) as Config;
        } else {
          return value as Config;
        }
      }
      const config = toConfig(tableData.value[i].Config);
      SubTitle.value = '编辑订阅'
      Subname.value = tableData.value[i].Name
      oldSubname.value = Subname.value
      if (config.udp)  {
        checkList.value.push('udp')
      }
      if (config.cert)  {
        checkList.value.push('cert')
      }
      Clash.value = config.clash
      Surge.value = config.surge
      dialogVisible.value = true
      value1.value = tableData.value[i].Nodes.map((item) => item.Name)
      sortSubID.value = tableData.value[i].ID
      sortValue.value = tableData.value[i].Nodes;
      sortValue.value.forEach((node, index) => {
        node.sort = index + 1; // 初始化 sort 值为当前顺序
      });
    }
  }
}

// 关闭对话框时清空数据
const handleEditDialogClose = () => {
  dialogVisible.value = false;
  sortSubID.value = null; // 清空订阅 ID
  sortValue.value = []; // 清空排序数据
  value1.value = []; // 清空已选节点
};


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
      getsubs()
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
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
const handleBase64 = (text: string) => {
  return  window.btoa(unescape(encodeURIComponent(text)));
}
const ClientDiaLog = ref(false)
const ClientList = ['v2ray','clash','surge'] // 客户端列表
const ClientUrls = ref<Record<string, string>>({})
const ClientUrl = ref('')
const handleClient = (name:string) => {
  let serverAddress = location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '');
  ClientDiaLog.value = true
  ClientUrl.value = `${serverAddress}/c/?token=${md5(name)}`
  ClientList.forEach((item:string) => {
    ClientUrls.value[item]=`${serverAddress}/c/?token=${md5(name)}`
  })
}

const Qrdialog = ref(false)
const QrTitle = ref('')
const handleQrcode = (url:string,title:string)=>{
  Qrdialog.value = true
  qrcode.value = url
  QrTitle.value = title
}
const OpenUrl = (url:string) => {
  window.open(url)
}
const clientradio = ref('1')

// 处理节点选择
const handleNodeSelection = () => {
  sortValue.value = NodesList.value.filter(node => value1.value.includes(node.Name));
};

const onDragEnd = () => {
  sortValue.value.forEach((node, index) => {
    node.sort = index + 1; // 更新排序值
  });
};

// 保存排序方法
const saveSortOrder = async () => {
  try {
    const nodesJson = JSON.stringify(sortValue.value.map((node) => ({
      nodeID: node.ID,
      sort: node.sort,
    })));
    await UpdateSubSort({
      subId: sortSubID.value,
      nodes: nodesJson,
    });
    ElMessage.success('排序更新成功');
  } catch (error) {
    console.error('Error updating nodes sort:', error);
    ElMessage.error('排序更新失败');
  }
  getsubs()
};

</script>

<template>
  <div>
    <el-dialog v-model="Qrdialog" width="300px" style="text-align: center" :title="QrTitle">
      <qrcode-vue :value="qrcode"  :size="200" level="H" />
      <el-input
      v-model="qrcode"
      >
      </el-input>
      <el-button @click="copyUrl(qrcode)">复制</el-button>
      <el-button @click="OpenUrl(qrcode)">打开</el-button>
    </el-dialog>

    <el-dialog v-model="ClientDiaLog" title="客户端(点击二维码获取地址)" style="text-align: center" >
      <el-row>
        <el-col>
        <el-tag type="success" size="large">自动识别</el-tag>
        <el-button @click="handleQrcode(ClientUrl,'自动识别客户端')">二维码</el-button>
      </el-col>
        <el-col v-for="(item,index) in ClientUrls" style="margin-bottom:10px;">
          <el-tag type="success" size="large">{{index}}</el-tag>
          <el-button @click="handleQrcode(`${item}&client=${index}`,index)">二维码</el-button>
        </el-col>
        </el-row>
    </el-dialog>

    <el-dialog v-model="iplogsdialog" title="访问记录" width="80%" draggable>
  <template #footer>
    <div class="dialog-footer">
      <el-table :data="IplogsList" border style="width: 100%">
        <el-table-column prop="IP" label="Ip" />
        <el-table-column prop="Count" label="总访问次数" />
        <el-table-column prop="Addr" label="来源" />
        <el-table-column prop="Date" label="最近时间" />
      </el-table>
    </div>
  </template>
</el-dialog>
    <el-dialog
    v-model="dialogVisible"
    :title="SubTitle"
  >
  <el-input v-model="Subname" placeholder="请输入订阅名称" />

  <el-row >
  <el-tag type="primary">clash模版选择</el-tag>
  <el-radio-group v-model="clientradio" class="ml-4">
      <el-radio value="1">本地</el-radio>
      <el-radio value="2">url链接</el-radio>
    </el-radio-group>
  <el-select v-model="Clash" placeholder="clash模版文件"  v-if="clientradio === '1'">
    <el-option v-for="template in templist" :key="template.file" :label="template.file" :value="'./template/'+template.file" />
  </el-select>
  <el-input v-model="Clash" placeholder="clash模版文件"  v-else />
</el-row>
<el-row >
  <el-tag type="primary">surge模版选择</el-tag>
  <el-radio-group v-model="clientradio" class="ml-4">
      <el-radio value="1">本地</el-radio>
      <el-radio value="2">url链接</el-radio>
    </el-radio-group>
  <el-select v-model="Surge" placeholder="surge模版文件"  v-if="clientradio === '1'">
    <el-option v-for="template in templist" :key="template.file" :label="template.file" :value="'./template/'+template.file" />
  </el-select>
  <el-input v-model="Surge" placeholder="surge模版文件"  v-else />
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
      @change="handleNodeSelection"
    >
      <el-option
        v-for="item in NodesList"
        :key="item.Name"
        :label="item.Name"
        :value="item.Name"
      />
    </el-select>
  </div>
  <!-- 添加节点排序 -->
  <div class="m-4">
    <p>节点排序</p>
    <!-- 使用 vuedraggable 实现拖拽排序 -->
    <draggable v-model="sortValue" tag="ul" item-key="ID" @end="onDragEnd">
      <template #item="{ element, index }">
        <li class="draggable-item">
          <!-- 使用 el-tag 展示排序数字 -->
          <el-tag type="info" class="sort-tag">
            {{ index + 1 }}
          </el-tag>
          <!-- 节点名称 -->
          <el-tag>{{ element.Name }}</el-tag>
        </li>
      </template>
    </draggable>
  </div>
    <template #footer>
      <div class="dialog-footer">
        <!-- 根据 subId 是否为空动态禁用按钮 -->
        <el-button type="primary" @click="saveSortOrder" :disabled="!sortSubID">保存排序</el-button>
        <el-button @click="handleEditDialogClose">关闭</el-button>
        <el-button type="primary" @click="addSubs">确定</el-button>
      </div>
    </template>
  </el-dialog>
    <el-card>
    <el-button type="primary" @click="handleAddSub">添加订阅</el-button>
    <div style="margin-bottom: 10px"></div>

      <el-table ref="table"
      :data="currentTableData"
      style="width: 100%"
      stripe
      @selection-change="handleSelectionChange"
      row-key="ID"
      :tree-props="{children: 'Nodes'}"
      >
    <el-table-column type="selection" fixed prop="ID" label="id"  />
    <el-table-column prop="Name" label="订阅名称 / 节点"  >
    <template #default="{row}">
      <el-tag :type="!row.Nodes ? 'success' : 'primary'" >{{row.Name}}</el-tag>
        </template>
    </el-table-column>
    <el-table-column prop="Link" label="链接" :show-overflow-tooltip="true" >
      <template #default="{row}">
        <div v-if="row.Nodes">
          <el-link type="primary" size="small" @click="handleClient(row.Name)">客户端</el-link>
        </div>
        </template>
      </el-table-column>

    <el-table-column prop="CreateDate" label="创建时间" sortable  />
    <el-table-column  label="操作" width="120">
      <template #default="scope">
        <div v-if="scope.row.Nodes">
          <el-button link type="primary" size="small" @click="handleIplogs(scope.row)">记录</el-button>
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
  <el-button link type="primary" size="small" @click="handleDel(scope.row)">删除</el-button>
        </div>
        <div v-else>
          <el-button link type="primary" size="small" @click="copyInfo(scope.row)">复制</el-button>
        </div>

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
