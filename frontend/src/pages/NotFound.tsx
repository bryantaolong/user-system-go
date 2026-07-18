import { useNavigate } from 'react-router-dom'
import { Button } from '@arco-design/web-react'

const NotFound = () => {
  const navigate = useNavigate()

  return (
    <div className="not-found-container">
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在。</p>
      <Button type="primary" onClick={() => navigate('/')}>
        返回首页
      </Button>
    </div>
  )
}

export default NotFound
